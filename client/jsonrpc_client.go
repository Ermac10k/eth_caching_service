package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/karlseguin/ccache/v2"
	"my.opera.eth.test/model"
)

const contentType = "application/json"

// JRClient is the object to request blocks from an ether node
type JRClient struct {
	url              string
	preformattedBody string
	cache            *ccache.Cache
	lastBlockNumber  *big.Int
	lock             sync.RWMutex
}

// NewJRClient is the JRClient constructor
func NewJRClient(url string, cache *ccache.Cache) (*JRClient, error) {
	n := new(big.Int)
	c := &JRClient{
		url:              url,
		preformattedBody: "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"%s\", true],\"id\":1}",
		cache:            cache,
		lastBlockNumber:  new(big.Int),
	}
	b, err := c.GetBlockBy("latest")
	if err != nil {
		return nil, err
	}
	_, ok := n.SetString(b.Number, 0)
	if !ok {
		return nil, fmt.Errorf("impossible to get latest block number")
	}
	c.lock.Lock()
	c.lastBlockNumber = c.lastBlockNumber.Set(n)
	c.lock.Unlock()
	return c, nil
}

// GetBlockBy is the GET method to request the latest block from eth chain
// identifier - can be a hex number in string format or the 'latest' tag
func (c *JRClient) GetBlockBy(identifier string) (*model.Block, error) {
	if identifier != "latest" {
		numID := new(big.Int)
		_, ok := numID.SetString(identifier, 0)
		if ok {

			c.lock.RLock()
			ln := numID.Set(c.lastBlockNumber)
			c.lock.RUnlock()

			cmp := new(big.Int).Sub(ln, numID).Cmp(big.NewInt(20))
			if cmp > 0 {
				log.Printf("check cache for block with number %s\n", identifier)
				cached := c.cache.Get(identifier)
				if cached != nil {
					log.Printf("block with number %s found in cache\n", identifier)
					return cached.Value().(*model.Block), nil
				}
				log.Printf("block with number %s not found in cache. requesting ethereum\n", identifier)
				b, err := c.receiveBlockStruct(identifier)
				if err != nil {
					return nil, err
				}

				// update cache concurrently
				go func() {
					log.Printf("update cache with block by number %s", identifier)
					c.cache.Set(identifier, b, time.Duration(math.MaxInt64))
				}()

				return b, nil
			}
		}
	}

	b, err := c.receiveBlockStruct(identifier)
	if err != nil {
		return nil, err
	}

	go c.updateLastNumber(b.Number)
	return b, nil

}

func (c *JRClient) updateLastNumber(new string) {
	log.Printf("Update latest block number with %s\n", new)
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lastBlockNumber.SetString(new, 0)
}

func (c *JRClient) receiveBlockStruct(identifier string) (*model.Block, error) {
	respBody, err := c.getBlockBytes(identifier)
	if err != nil {
		return nil, err
	}
	resp, err := c.bytesToBlockJSON(respBody, identifier)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, &model.ResponseContentError{
			Message: resp.Error.Message,
		}
	}
	if resp.Result == nil {
		return nil, &model.ResponseContentError{
			Message: "a resulting block in a response is empty because of unknown reason",
		}
	}
	return resp.Result, nil
}

func (c *JRClient) getBlockBytes(param string) ([]byte, error) {
	data := strings.NewReader(fmt.Sprintf(c.preformattedBody, param))
	log.Printf("request for a block by identifier %s\n", param)
	resp, err := http.Post(c.url, contentType, data)
	if err != nil {
		log.Printf(
			"an error (%s) occured while requesting a block by identifier %s\n",
			err.Error(),
			param,
		)
		return nil, err
	}
	log.Printf(
		"received answer for a block with identifier %s, Status code = %d\n",
		param,
		resp.StatusCode,
	)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"an error (%s) occured while reading a response with a block by identifier %s\n",
			err.Error(),
			param,
		)
		return nil, err
	}
	return body, nil
}

func (c *JRClient) bytesToBlockJSON(data []byte, param string) (*model.RespJSON, error) {
	resp := new(model.RespJSON)
	if err := json.Unmarshal(data, resp); err != nil {
		log.Printf(
			"an error (%s) occured while creating json of a block to  by identifier %s\n",
			err.Error(),
			param,
		)
		return nil, err
	}
	return resp, nil
}

// GetTransactionByHash finds a particular transaction in a requested block
func (c *JRClient) GetTransactionByHash(block *model.Block, hash string) (*model.Transaction, error) {
	log.Printf(
		"searching in a block with a number %s for a transaction with hash %s\n",
		block.Number,
		hash,
	)
	for _, t := range block.Transactions {
		if t.Hash == hash {
			log.Printf(
				"the transaction with the hash %s found in the block %s",
				hash,
				block.Number,
			)
			return t, nil
		}
	}
	log.Printf(
		"the transaction with the hash %s not found in the block %s",
		hash,
		block.Number,
	)
	return nil, &model.NotFoundHashTransactionError{BlockHash: block.Hash, Hash: hash}
}

// GetTransactionByIndex finds a particular transaction in a requested block
func (c *JRClient) GetTransactionByIndex(block *model.Block, index uint64) (*model.Transaction, error) {
	log.Printf(
		"searching in a block with a number %s for a transaction with index %d\n",
		block.Number,
		index,
	)
	indexHex := fmt.Sprintf("0x%x", index)
	// Nowhere or nothing to find
	if block == nil || block.Transactions == nil ||
		len(block.Transactions) == 0 ||
		uint64(len(block.Transactions)) <= index {
		log.Printf(
			"the transaction with the index %d not found in the block %s",
			index,
			block.Number,
		)
		return nil, &model.NotFoundIDTransactionError{BlockHash: block.Hash, ID: fmt.Sprintf("0x%d", index)}
	}
	for _, t := range block.Transactions {
		if t.TransactionIndex == indexHex {
			return t, nil
		}
	}
	log.Printf(
		"the transaction with the index %d not found in the block %s",
		index,
		block.Number,
	)
	return nil, &model.NotFoundIDTransactionError{BlockHash: block.Hash, ID: fmt.Sprintf("0x%d", index)}
}
