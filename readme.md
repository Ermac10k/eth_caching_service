# Caching proxy for eth_getBlockByNumber

## Endpoints

supports next endpoints:
+ `/block/latest` - GET a latest block in a chain
+ `/block/{number}` - GET a block with filed "number"={number}, where number is decimal
+ `/block/latest/txs/{identifierT}` - GET a transaction from a latest block by it's "hash" or "transactionIndex" field. So "hash" is string of "0x..." format and "transactionIndex" is decimal
+ `/block/{number}/txs/{identifier}` - GET a transaction from a block with filed "number"={number} by it's "hash" or "transactionIndex" field. So "hash" is string of "0x..." format and "transactionIndex" is decimal

## Run Args

*All flags are optional*
+ `-host` - a hostname to start a service. **default**=`localhost`
+ `-port` - "a port to start service. **default**=`8080`
+ `-node` - "an address of an ether node to request blocks. **default**=`https://cloudflare-eth.com`
+ `-csize` - "a cache size to store blocks. **default is** `MaxInt64`

## Techstack

+ **github.com/valyala/fasthttp** - as an HTTP server. Because it's fast
+ **github.com/fasthttp/router** - as a router over fasthttp to handle endpoints. Becouse it's fast and handy
+ **github.com/karlseguin/ccache/v2** - as a LRU cache. Because it's handy, reliable and it's possible to tune sizing. *I don't implement the method Size because there wasn't a place for it in this task. But it's a lucky find to me*  
I've made some experiments and ensured that it has a good control over memory overheads and concurrency races. Also it's being suported till today. It has a few issues on Github

---
P.S. Some requirements were set via e-mail.  