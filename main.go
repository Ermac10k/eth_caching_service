package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/karlseguin/ccache/v2"

	"my.eth.test/client"
	"my.eth.test/logger"
	"my.eth.test/server"
)

func main() {
	host := flag.String("host", "localhost", "a hostname to start a service. default=localhost")
	port := flag.Uint("port", 8080, "a port to start service. default=8080")
	etherAddr := flag.String("node", "https://cloudflare-eth.com", "an address of an ether node to request blocks. default=https://cloudflare-eth.com")
	cacheSize := flag.Int64("csize", 0, "a cache size to store blocks. default=MaxInt64")
	flag.Parse()

	// create cache
	var size int64
	if *cacheSize > 0 {
		size = *cacheSize
	} else {
		size = math.MaxInt64
	}
	cache := ccache.New(ccache.Configure().Buckets(256).ItemsToPrune(100).MaxSize(size))

	// create client to request blocks
	log.SetFlags(0)
	log.SetOutput(new(logger.Logger))
	locclient, err := client.NewJRClient(*etherAddr, cache)
	if err != nil {
		log.Fatal(err)
	}

	// create server
	server := server.NewRouterToServe(*host, fmt.Sprint(*port), locclient)
	log.Fatal(server.Serve())
}
