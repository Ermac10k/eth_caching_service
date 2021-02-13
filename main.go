package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"my.opera.eth.test/client"
)

func main() {
	// ch := make(chan string)
	// go func() {
	// 	data := strings.NewReader("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"latest\", true],\"id\":1}")
	// 	resp, _ := http.Post("https://cloudflare-eth.com", "application/json", data)
	// 	body, _ := ioutil.ReadAll(resp.Body)
	// 	ch <- string(body)
	// }()
	// fmt.Println(<-ch)
	// close(ch)
	// fmt.Println()

	// create client to request blocks
	locclient := client.NewJRClient("https://cloudflare-eth.com")

	index := "0x2244" // 8772
	num := new(big.Int)
	if _, ok := num.SetString(index, 0); !ok {
		log.Fatal("0x2244")
	}

	b, err := locclient.GetBlockBy("latest")
	if err != nil {
		log.Fatal(err)
	}

	// t, _ := locclient.GetTransactionByHash(b, b.Transactions[0].Hash)
	t, _ := locclient.GetTransactionByIndex(b, 3)

	// bb, _ := json.Marshal(b)
	tt, _ := json.Marshal(t)
	// fmt.Println(string(bb))
	
	fmt.Println(string(tt))
}
