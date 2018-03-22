package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

func handleCon(conn net.Conn) {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	buf := make([]byte, 0)
	tmp := make([]byte, 256)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)
	}
	log.Printf("%s", string(buf))
	//sleep randomly
	time.Sleep(time.Duration(r.Intn(50)) * time.Millisecond)
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
