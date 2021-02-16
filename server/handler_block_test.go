package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/karlseguin/ccache/v2"
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

var testCasesBlocks = map[string]string{
	"TestBlockIDOne":      "/block/1",
	"TestBlockIDLatest":   "/block/latest",
	"TestBlockIDNegative": "/block/-1",
	"TestBlockIDString":   "/block/ff", // ff is not only a string but a hexadecimal number either. double check.
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

func TestIncorrectEtherAddress(t *testing.T) {
	cache := ccache.New(ccache.Configure().Buckets(8).ItemsToPrune(1).MaxSize(1))
	_, err := client.NewJRClient("https://cloudflare-eth.c", cache)
	if err == nil {
		t.Error(fmt.Errorf("expected error"))
	}
}

func commonPart(t *testing.T, tcs map[string]string) (*http.Response, []byte) {
	cache := ccache.New(ccache.Configure().Buckets(8).ItemsToPrune(1).MaxSize(1))
	var c, err = client.NewJRClient("https://cloudflare-eth.com", cache)
	if err != nil {
		t.Error(err)
		return nil, nil
	}
	s := NewRouterToServe("test", "", c)
	r, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://%s:%s%s", s.host, s.port, tcs[t.Name()]),
		nil,
	)
	res, err := serve(RegisterHandler(s), r)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return nil, nil
	}

	return res, body
}
