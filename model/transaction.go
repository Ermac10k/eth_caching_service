package model

// Transaction json response
type Transaction struct {
	BlockHash string `josn:"blockHash"`
	BlockNumber string `json:"blockNumber"`
	From string `json:"from"`
	Gas string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Hash string `json:"hash"`
	Input string `json:"input"`
	Nonce string `json:"nonce"`
	To string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value string `json:"value"`
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`
}