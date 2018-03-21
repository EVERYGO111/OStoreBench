package main

import (
	"driver"
	"log"
	"net"
	"workloads"
)

const (
	// MAX_FILE_COUNT is the max count of files
	MAX_FILE_COUNT = 1048576
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
	//r := rand.New(rand.NewSource(int64(time.Now().Second())))
	//zipf := rand.NewZipf(r, 2.7, 1048576/5, 1048576)
	//// data := make([]int, 0)
	//for i := 0; i != MAX_FILE_COUNT; i++ {
	//	item := int(zipf.Uint64())
	//	fmt.Println(item)
	//	//data = append(data, item)
	//}

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
	//client := driver.NewHTTPClient("localhost", 8000)
	//w := workloads.NewPoissonWorkloads(50, 10485760, MAX_FILE_COUNT, client)
	//w.Start()

	type WorkloadConfig struct {
		rwRate       float64 //R/W rate
		totalProcess int64   //total goroutine

		driver       driver.Driver
		fileSizeType workloads.DistributionType
		iatType      workloads.DistributionType
		requestType  workloads.DistributionType

		requestNum int64
	}
	workloadConf := workloads.WorkloadConfig{
		WriteRate:       0.3,
		TotalProcess: 10,
		Driver:       driver.NewFakeDriver("localhost", 8000),
		FileSizeType: workloads.ZIPF,
		IatType:      workloads.NEGATIVE_EXP,
		RequestType:  workloads.LOGNORMAL,
		RequestNum:   10000,
	}
	workloads := workloads.NewWorkload(workloadConf)
	workloads.Start()
}
