package model

import "fmt"

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

// InvalidIdentifierError to report that a block identifier to request Block from the ether node is invalid 
type InvalidIdentifierError struct {
	Identifier string
}

func (err *InvalidIdentifierError) Error() string {
	return fmt.Sprintf("an identifier: '%s' is invalid", err.Identifier)
}

// ResponseContentError to report that there is an error in an ethereum node response
type ResponseContentError struct {
	Message string
}

func (err *ResponseContentError) Error() string {
	return fmt.Sprintf("Ethereum node has returned an error with message: %s", err.Message)
}