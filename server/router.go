package server

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"my.opera.eth.test/client"
)

// RouterToServe is the service object
type RouterToServe struct {
	host   string
	port   string
	client *client.JRClient
}

// NewRouterToServe is the constructor of the RoterToServe obj
func NewRouterToServe(hostname string, port string, c *client.JRClient) *RouterToServe {
	return &RouterToServe{
		hostname,
		port,
		c,
	}
}

// Serve registers handlers and starts the service
func (s *RouterToServe) Serve() error {
	initAddr := fmt.Sprintf("%s:%s", s.host, s.port)
	handler := RegisterHandler(s)
	return fasthttp.ListenAndServe(initAddr, handler)
}

// RegisterHandler registers new router and returns a handler from it
func RegisterHandler(s *RouterToServe) func(*fasthttp.RequestCtx) {
	r := router.New()
	r.GET("/block/{identifier}", s.requestBlock)
	r.GET("/block/{identifierB}/txs/{identifierT}", s.requestBlockAndFindTransaction)
	return r.Handler
}
