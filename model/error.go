package model

import "fmt"
// Error json
type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}


// NotFoundIDTransactionError to report that requested transaction not found in a particular block
type NotFoundIDTransactionError struct {
	BlockHash string
	ID string
}

func (err *NotFoundIDTransactionError) Error() string {
	return fmt.Sprintf("the transaction with ID=%s not found in a requested block (blockHash=%s)", err.ID, err.BlockHash)
}


// NotFoundHashTransactionError to report that requested transaction not found in a particular block
type NotFoundHashTransactionError struct {
	BlockHash string
	Hash string
}

func (err *NotFoundHashTransactionError) Error() string {
	return fmt.Sprintf("the transaction with Hash=%s not found in a requested block (blockHash=%s)", err.Hash, err.BlockHash)
}

// InvalidBlockIdentifierError to report that a block identifier to request Block from the ether node is invalid 
type InvalidBlockIdentifierError struct {
	Identifier string
}

func (err *InvalidBlockIdentifierError) Error() string {
	return fmt.Sprintf(" a block identifier: '%s' to request Block from the ether node is invalid", err.Identifier)
}