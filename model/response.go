package model

// RespJSON is the dto to unmarshal json resp
type RespJSON struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  *Block    `json:"result"`
	ID      int       `json:"id"`
	Error   *EthError `json:"error"`
}

// EthError json
type EthError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Block json response
// I am strictly follow the Cloudflare Docs output format
// so this is a combination of the Header and the Block structs
type Block struct {
	NoTransactionBlock
	Transactions []*Transaction `json:"transactions"`
}

// ToShowcase is the converter from a whole block to a block with Transactions array
// that contains transactions' hashes only
func (b *Block) ToShowcase() *ShowcaseBlock {
	// sort transactions to ensure by Number. This requirement was introduced in an additional email from Igor
	txs := hashes(b.Transactions, b.Hash)
	return &ShowcaseBlock{NoTransactionBlock: b.NoTransactionBlock, Transactions: txs}
}

func hashes(txs []*Transaction, blockHash string) []string {
	tHashes := make([]string, len(txs))
	for i, t := range txs{
		tHashes[i] = t.Hash
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
