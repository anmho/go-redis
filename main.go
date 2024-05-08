package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const port = 6379

func main() {

	addrString := fmt.Sprintf(":%d", port)

	l, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("listening on port %d\n", port)

	conn, err := l.Accept()
	fmt.Println("Accepted connection from client")

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		fmt.Println(string(buf))
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln("error reading from client", err.Error())
		}
		conn.Write([]byte("+OK\r\n"))
	}

}
