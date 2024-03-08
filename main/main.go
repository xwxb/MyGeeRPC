package main

import (
	"encoding/json"
	"fmt"
	geerpc "github.com/xwxb/MyGeeRPC"
	"github.com/xwxb/MyGeeRPC/codec"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String() // 信道传输成功连接后的地址，实现了并发同步下的确定执行瞬时
	geerpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple geerpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)

	// 最过拟合难理解的点是，conn 本身就实现了 writer 是直接往里写的，这里不是抽象成了显示请求
	// send options
	// json encoder convert struct to json format and write to conn
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption) // return a encoder write to conn, and encode the DefaultOption
	cc := codec.NewGobCodec(conn)                          // client side codec
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}

		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)

		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
