package model

import (
	"log"
	"sort"
	"strconv"
)

// RespJSON is the dto to unmarshal json resp
type RespJSON struct {
	JSONRPC string `json:"jsonrpc"`
	Result  *Block `json:"result"`
	ID      int    `json:"id"`
}

// Block json response
// I am strictly follow the Cloudflare Docs output format
// so this is a combination of the Header and the Block structs
type Block struct {
	NoTransactionBlock
	Transactions []*Transaction `json:"transactions"`
}

type txMetaPair struct {
	Index uint64
	Hash  string
}

// ToShowcase is the converter from a whole block to a block with Transactions array
// that contains transactions' hashes only
func (b *Block) ToShowcase() *ShowcaseBlock {
	// sort transactions to ensure by Number. This requirement was introduced in an additional email from Igor
	txs := sortHashes(b.Transactions, b.Hash)
	return &ShowcaseBlock{NoTransactionBlock: b.NoTransactionBlock, Transactions: txs}
}

func sortHashes(txs []*Transaction, blockHash string) []string {
	var errDecision error
	tHashes := make([]string, len(txs))
	pairs := make([]txMetaPair, len(txs))
	for i, t := range txs {
		u, err := strconv.ParseUint(t.TransactionIndex, 16, 64)
		if err != nil {
			errDecision = err
			break
		}
		pairs[i] = txMetaPair{u, t.Hash}
	}

	if errDecision == nil {
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Index < pairs[j].Index
		})
		for i, p := range pairs {
			tHashes[i] = p.Hash
		}
	} else {
		log.Printf("the transactions array of block %s is unsortable\nerror:%s\n", blockHash, errDecision.Error())
		for i, t := range txs {
			tHashes[i] = t.Hash
		}
	}

	return tHashes
}

// NoTransactionBlock is the dto to construct a got result and a showing one
type NoTransactionBlock struct {
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}
