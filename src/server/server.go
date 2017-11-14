package main

import (
	"log"
	"net"
)

func handleCon(conn net.Conn) {
	log.Printf("%s", conn.RemoteAddr().String())
	defer conn.Close()
}
func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleCon(conn)
	}
}
