package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"my.opera.eth.test/client"
	"my.opera.eth.test/model"
)

func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

var c = client.NewJRClient("https://cloudflare-eth.com")
var s = NewRouterToServe("test", "", c)
var testCasesBlocks = map[string]*http.Request{
	"TestBlockIDOne": func(str string) *http.Request {
		r, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s%s", s.host, s.port, str), nil)
		return r
	}("/block/1"),
	"TestBlockIDLatest": func(str string) *http.Request {
		r, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s%s", s.host, s.port, str), nil)
		return r
	}("/block/latest"),
	"TestBlockIDNegative": func(str string) *http.Request {
		r, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s%s", s.host, s.port, str), nil)
		return r
	}("/block/-1"),
	"TestBlockIDString": func(str string) *http.Request {
		r, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s%s", s.host, s.port, str), nil)
		return r
	}("/block/ff"), // ff is not only a string but a hexadecimal number either. double check.
}

func TestBlockIDOne(t *testing.T) {
	_, body := commonPart(t, testCasesBlocks)
	b := new(model.ShowcaseBlock)
	err := json.Unmarshal(body, b)
	if err != nil {
		t.Error(err)
	}

	if b.Hash == "" {
		t.Error("no response!")
	}
}

func TestBlockIDLatest(t *testing.T) {
	_, body := commonPart(t, testCasesBlocks)
	b := new(model.ShowcaseBlock)

	err := json.Unmarshal(body, b)
	if err != nil {
		t.Error(err)
	}

	if b.Hash == "" {
		t.Error("no response!")
	}
}

func TestBlockIDNegative(t *testing.T) {
	resp, body := commonPart(t, testCasesBlocks)
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

func TestBlockIDString(t *testing.T) {
	resp, body := commonPart(t, testCasesBlocks)
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

func commonPart(t *testing.T, tcs map[string]*http.Request) (*http.Response, []byte) {
	r := tcs[t.Name()]

	res, err := serve(RegisterHandler(s), r)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	return res, body
}
