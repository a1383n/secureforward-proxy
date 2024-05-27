package main

import (
	"fmt"
	"log"
	"net"
	"secureforward-proxy/src"
)

func main() {
	l, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on :443")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println("Connection accepted from ", conn.RemoteAddr())
		go src.HandleConnection(conn)
	}
}
