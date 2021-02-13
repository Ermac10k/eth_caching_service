package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"

	"my.opera.eth.test/model"
)

const contentType = "application/json"

// jrClient is the object to request blocks from an ether node
type jrClient struct {
	url              string
	preformattedBody string
}

// NewJRClient is the JRClient constructor
func NewJRClient(url string) *jrClient {
	return &jrClient{url: url, preformattedBody: "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"%s\", true],\"id\":1}"}
}

// GetBlockBy is the GET method to request the latest block from eth chain
// identifier - can be a hex number in string format or the 'latest' tag
func (c *jrClient) GetBlockBy(identifier string) (*model.Block, error) {
	// validate the identifier value
	if identifier != "latest" {
		if _, ok := (new(big.Int)).SetString(identifier, 0); !ok {
			return nil, &model.InvalidBlockIdentifierError{Identifier: identifier}
		}
	}

	respBody, err := c.getBlockBytes(identifier)
	if err != nil {
		return nil, err
	}

	resp, err := c.bytesToBlockJSON(respBody)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, &model.ResponseContentError{Message: resp.Error.Message}
	}
	if resp.Result == nil {
		return nil, &model.ResponseContentError{Message: "a resulting block in a response is empty because of unknown reason"}
	}
	return resp.Result, nil
}

func (c *jrClient) getBlockBytes(param string) ([]byte, error) {
	data := strings.NewReader(fmt.Sprintf(c.preformattedBody, param))
	resp, err := http.Post(c.url, contentType, data)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *jrClient) bytesToBlockJSON(data []byte) (*model.RespJSON, error) {
	resp := new(model.RespJSON)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTransactionByHash finds a particular transaction in a requested block
func (c *jrClient) GetTransactionByHash(block *model.Block, hash string) (*model.Transaction, error) {
	log.Printf("request=%s\n", hash)
	for _, t := range block.Transactions {
		log.Printf("iterate=%s\n", t.Hash)
		if t.Hash == hash {
			return t, nil
		}
	}
	return nil, &model.NotFoundHashTransactionError{BlockHash: block.Hash, Hash: hash}
}

// GetTransactionByIndex finds a particular transaction in a requested block
func (c *jrClient) GetTransactionByIndex(block *model.Block, index uint64) (*model.Transaction, error) {
	indexHex := fmt.Sprintf("0x%x", index)
	// Nowhere or nothing to find
	if block == nil || block.Transactions == nil ||
		len(block.Transactions) == 0 ||
		int64(len(block.Transactions)) <= int64(index) {
		return nil, &model.NotFoundIDTransactionError{BlockHash: block.Hash, ID: fmt.Sprintf("0x%d", index)}
	}
	for _, t := range block.Transactions {
		if t.TransactionIndex == indexHex {
			return t, nil
		}
	}

	return nil, &model.NotFoundIDTransactionError{BlockHash: block.Hash, ID: fmt.Sprintf("0x%d", index)}
}
