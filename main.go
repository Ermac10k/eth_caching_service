package main

import (
	"log"

	"my.opera.eth.test/client"
	"my.opera.eth.test/server"
)

func main() {
	// 8772 //0xb4e573 //11855219 //0x29489800f624b64b975af75bde520c5a70a21848920b5483d463c25c3b22ac0b
	// num := new(big.Int)
	// num, _ = num.SetString("0xb4e573", 0)
	// fmt.Println(num.Int64())

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

	// create server
	server := server.NewRouterToServe("localhost", "8080", locclient)
	log.Fatal(server.Serve())
}
