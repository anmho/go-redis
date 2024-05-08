package main

import (
	"fmt"
	"github.com/amho/go-redis/server"
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
		//buf := make([]byte, 1024)
		resp := server.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(value)
		//_, err := conn.Read(buf)
		//log.Println(string(buf))
		//resp.Read()

		//if err != nil {
		//	if err == io.EOF {
		//		break
		//	}
		//	log.Fatalln("error reading from client", err.Error())
		//}

		conn.Write([]byte("+OK\r\n"))
	}
}
