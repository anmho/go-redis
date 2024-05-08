package main

import (
	"fmt"
	"github.com/amho/go-redis/redis"
	"log"
	"net"
)

const port = 6379

func main() {

	addrString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("listening on port %d\n", port)
	conn, err := listener.Accept()
	fmt.Println("Accepted connection from client")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	for {
		resp := redis.NewResp(conn)
		req, err := resp.Read()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(req)

		writer := redis.NewWriter(conn)
		v := redis.NewString("OK")
		err = writer.Write(v)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
