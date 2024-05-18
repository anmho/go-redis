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

func serve(conn net.Conn, aof *Aof) {
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
			writer.Write(NewError(err.Error()))
			continue
		}
		command := req.array[0].bulk
		if command == "SET" || command == "HSET" {
			err := aof.Write(req)
			if err != nil {
				panic(err)
			}
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

	aof, err := NewAof("database.aof")
	if err != nil {
		panic(err)
	}
	defer aof.Close()
	err = aof.Read()
	if err != nil {
		panic(err)
	}

	// should prevent busy spinning
	for {
		conn, err := listener.Accept()
		fmt.Println("Accepted connection from client")
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			serve(conn, aof)
			defer conn.Close()
		}()
	}
}
