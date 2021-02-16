package server

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
	"my.opera.eth.test/model"
)

var testCasesTxs = map[string]string{
	"TestBlockIDNegativeT": "/block/-1/txs/1",
	"TestBlockIDStringT":   "/block/ff/txs/1", // ff is not only a string but a hexadecimal number either. double check.
	"TestTxIDNegative":     "/block/11855219/txs/-1",
	"TestTxIDString":       "/block/11855219/txs/ff", // ff is not only a string but a hexadecimal number either. double check.
	"TestTxByHash":         "/block/11855219/txs/0x29489800f624b64b975af75bde520c5a70a21848920b5483d463c25c3b22ac0b",
	"TestTxById":           "/block/11855219/txs/1",
}

func TestTxByHash(t *testing.T) {
	_, body := commonPart(t)
	tx := new(model.Transaction)
	err := json.Unmarshal(body, tx)
	if err != nil {
		t.Error(err)
	}

	if tx.TransactionIndex != "0x0" {
		t.Error("invalid response!")
	}
}

func TestTxById(t *testing.T) {
	_, body := commonPart(t)
	tx := new(model.Transaction)
	err := json.Unmarshal(body, tx)
	if err != nil {
		t.Error(err)
	}

	if tx.TransactionIndex != "0x1" {
		t.Error("invalid response!")
	}
}

func TestBlockIDNegativeT(t *testing.T) {
	resp, body := commonPart(t)
	if resp.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf(
			"Invalid status code: %d\nexpectd: %d",
			resp.StatusCode,
			fasthttp.StatusBadRequest,
		)
	}
	if string(body) != fmt.Sprintf(
		"an identifier: '%d' is invalid",
		-1,
	) {
		t.Errorf(
			"Invalid message: %s\nexpectd: %s",
			string(body),
			fmt.Sprintf("an identifier: '%d'  is invalid", -1),
		)
	}
}

func TestBlockIDStringT(t *testing.T) {
	resp, body := commonPart(t)
	if resp.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf(
			"Invalid status code: %d\nexpectd: %d",
			resp.StatusCode,
			fasthttp.StatusBadRequest,
		)
	}
	if string(body) != fmt.Sprintf(
		"an identifier: '%s' is invalid",
		"ff",
	) {
		t.Errorf(
			"Invalid message: %s\nexpectd: %s",
			string(body),
			fmt.Sprintf(
				"an identifier: '%s' is invalid",
				"ff",
			),
		)
	}
}

func TestTxIDNegative(t *testing.T) {
	resp, body := commonPart(t)
	if resp.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf(
			"Invalid status code: %d\nexpectd: %d",
			resp.StatusCode,
			fasthttp.StatusBadRequest,
		)
	}
	if string(body) != fmt.Sprintf(
		"an identifier: '%d' is invalid",
		-1,
	) {
		t.Errorf(
			"Invalid message: %s\nexpectd: %s",
			string(body),
			fmt.Sprintf("an identifier: '%d'  is invalid", -1),
		)
	}
}

func TestTxIDString(t *testing.T) {
	resp, body := commonPart(t)
	if resp.StatusCode != fasthttp.StatusBadRequest {
		t.Errorf(
			"Invalid status code: %d\nexpectd: %d",
			resp.StatusCode,
			fasthttp.StatusBadRequest,
		)
	}
	if string(body) != fmt.Sprintf(
		"an identifier: '%s' is invalid",
		"ff",
	) {
		t.Errorf(
			"Invalid message: %s\nexpectd: %s",
			string(body),
			fmt.Sprintf(
				"an identifier: '%s' is invalid",
				"ff",
			),
		)
	}
}
