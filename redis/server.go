package redis

import (
	"fmt"
	"log"
	"net"
)

type Redis struct {
}

func New() *Redis {
	return &Redis{}
}

func serve(conn net.Conn) {
	defer conn.Close()

	for {
		resp := NewResp(conn)
		req, err := resp.Read()
		if err != nil {
			log.Fatalln(err)
		}

		writer := NewWriter(conn)

		log.Println(req)

		result, err := Handle(req)
		if err != nil {
			log.Println(err)
			continue
		}
		err = writer.Write(result)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (r *Redis) ListenAndServe(port int) {
	addrString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		fmt.Println("Accepted connection from client")
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			serve(conn)
			defer conn.Close()
		}()
	}
}
