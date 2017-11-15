package main

import (
	"driver"
	"log"
	"net"
	"workloads"
)

const (
	// MAX_FILE_COUNT is the max count of files
	MAX_FILE_COUNT = 100000
)

func sendRequest() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte("hello world"))
	conn.Close()
}

func main() {
	// r := rand.New(rand.NewSource(int64(time.Now().Second())))
	// zipf := rand.NewZipf(r, 2.7, 25, 300)
	// data := make([]int, 0)
	// N := 120
	// for i := 0; i != N; i++ {
	// 	item := int(zipf.Uint64())
	// 	data = append(data, item)
	// }

	// p := distribution.NewNegativeExp(1000 / 50)
	// data := make([]int64, 0)
	// for i := 0; i < 8000; i++ {
	// 	item := p.Float64()
	// 	data = append(data, int64(item)*1000)
	// }
	//
	// for i := 0; i < len(data); i++ {
	// 	// fmt.Println(data[i])
	// 	sendRequest()
	// 	time.Sleep(time.Duration(data[i]) * time.Microsecond)
	// }
	client := driver.NewHTTPClient("localhost", 8000)
	w := workloads.NewPoissonWorkloads(50, 10485760, MAX_FILE_COUNT, client)
	w.Start()
}
