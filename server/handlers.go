package server

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
	"my.eth.test/model"
)

// GET /block/{identifier}
func (s *RouterToServe) requestBlock(ctx *fasthttp.RequestCtx) {
	identifier := ctx.UserValue("identifier").(string)
	if identifier != "latest" { // separates the 'latest' tag from numeric values
		if err := validateParam(&identifier); err != nil { // changes identifier format to 0x...
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			return
		}
	}
	block, err := s.client.GetBlockBy(identifier)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(block.ToShowcase())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.WriteString(string(resp))
}

func validateParam(identifier *string) error {
	if num, err := strconv.ParseUint(*identifier, 10, 64); err == nil {
		*identifier = fmt.Sprintf("0x%x", num)
	} else {
		return &model.InvalidIdentifierError{Identifier: *identifier}
	}
	return nil
}

// GET /block/{identifierB}/txs/{identifierT}
func (s *RouterToServe) requestBlockAndFindTransaction(ctx *fasthttp.RequestCtx) {
	var numericIDT uint64
	isHash := true // to separate an identifierT hash value from an identifierT numeric value
	idB := ctx.UserValue("identifierB").(string)
	if idB != "latest" { // separates the 'latest' tag from numeric values
		if err := validateParam(&idB); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			return
		}
	}
	idT := ctx.UserValue("identifierT").(string)
	if len(idT) < 2 || idT[:2] != "0x" { // separates hash values from numeric values
		num, err := strconv.ParseUint(idT, 10, 64)
		if err != nil {
			ctx.Error(
				(&model.InvalidIdentifierError{Identifier: idT}).Error(),
				fasthttp.StatusBadRequest,
			)
			return
		}
		isHash = false
		numericIDT = num
	}
	block, err := s.client.GetBlockBy(idB)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	var t *model.Transaction
	if isHash {
		t, err = s.client.GetTransactionByHash(block, idT)
	} else {
		t, err = s.client.GetTransactionByIndex(block, numericIDT)
	}

	resp, err := json.Marshal(t)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.WriteString(string(resp))
}
