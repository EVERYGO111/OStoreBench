package main

import (
	"driver"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"workloads"
)

const (
	// MAX_FILE_COUNT is the max count of files
	MAX_FILE_COUNT  = 1048576
	ONLINE_SERVICE  = "OnlineService"
	ARCHIVE_SERVICE = "Archive"
)

type BenchmarkOptions struct {
	target      *string
	concurrency *int

	//for haystack
	server *string

	//for swift
	authUrl  *string
	username *string
	apiKey   *string
	domain   *string
	project  *string

	workloadType *string
	filesize     *int
	requestNum   *int64 //for onlineservice

}

var (
	b BenchmarkOptions
)

func init() {
	b.target = flag.String("target", "weed", "weed, swift or ceph")
	b.server = flag.String("server", "localhost:9333", "server ip")
	b.concurrency = flag.Int("c", 16, "number of cocurrent read and write process")
	b.authUrl = flag.String("authUrl", "http://localhost:5000/v3", "auth url for swift")
	b.username = flag.String("username", "demo", "username for swift")
	b.apiKey = flag.String("apikey", "123456", "pasword or api key issued by the provider")
	b.domain = flag.String("domain", "default", "domain for swift")
	b.project = flag.String("project", "demo", "project name for swift")
	b.workloadType = flag.String("t", ONLINE_SERVICE, "1: OnlineService, 2: Archive ")
	b.filesize = flag.Int("filesize", 10240, "maximum file size")
	b.requestNum = flag.Int64("reqnum", 10000, "request number")
	flag.Parse()
}
func benchOnlineservice(driver driver.Driver) {
	//分布可调负载
	workloadConf := workloads.WorkloadConfig{
		WriteRate:    0.3,
		TotalProcess: int64(*b.concurrency),
		//Driver:       driver.NewWeeDriver("localhost:9333"),
		Driver:       driver,
		FileSizeType: workloads.ZIPF,
		IatType:      workloads.NEGATIVE_EXP, //the request rate will be lognormal distribution
		RequestType:  workloads.LOGNORMAL,
		RequestNum:   *b.requestNum,
	}
	workloads := workloads.NewWorkload(workloadConf)
	workloads.Start()
}

func benchArchiveService(driver driver.Driver) {
	//归档负载
	archWorkloadConf := workloads.ArchiveWorkloadConfig{
		RequestGroups:    100,
		GroupInteralTime: 5,
		MinFileSize:      10,
		MaxFileSize:      1024,
		MaxFilesPerGroup: 10,
		ProcessNum:       *b.concurrency,
	}
	//workload := workloads.NewArchWorkload(driver.NewFakeDriver("localhost", 8000), &archWorkloadConf)
	workload := workloads.NewArchWorkload(driver, &archWorkloadConf)
	workload.Start()
}

func loadDriverByType(t string) (driver.Driver, error) {
	switch t {
	case "weed":
		return driver.NewWeeDriver(*b.server), nil
	case "swift":
		return driver.NewSwiftDriver(*b.username, *b.apiKey, *b.authUrl, *b.domain, *b.project), nil
	case "ceph":
		return driver.NewCephDriver(*b.username, *b.apiKey, *b.authUrl), nil
	default:
		return nil, fmt.Errorf("Ceph driver has not been implemented yet!")
	}
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

	//分布可调负载
	//workloadConf := workloads.WorkloadConfig{
	//	WriteRate:    0.3,
	//	TotalProcess: 10,
	//	//Driver:       driver.NewWeeDriver("localhost:9333"),
	//	Driver:       driver.NewSwiftDriver("demo","123456","http://172.16.1.92:5000/v3","default","demo"),
	//	FileSizeType: workloads.ZIPF,
	//	IatType:      workloads.NEGATIVE_EXP, //the request rate will be lognormal distribution
	//	RequestType:  workloads.LOGNORMAL,
	//	RequestNum:   10000,
	//}
	//workloads := workloads.NewWorkload(workloadConf)
	//workloads.Start()

	//归档负载
	//archWorkloadConf := workloads.ArchiveWorkloadConfig{
	//	RequestGroups:    100,
	//	GroupInteralTime: 5,
	//	MinFileSize:      10,
	//	MaxFileSize:      1024,
	//	MaxFilesPerGroup: 10,
	//	ProcessNum:       16,
	//}
	////workload := workloads.NewArchWorkload(driver.NewFakeDriver("localhost", 8000), &archWorkloadConf)
	//workload := workloads.NewArchWorkload(driver.NewSwiftDriver("demo", "123456", "http://172.16.1.92:5000/v3", "default", "demo"), &archWorkloadConf)
	//workload.Start()
	//
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	driver, err := loadDriverByType(*b.target)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	switch *b.workloadType {
	case ONLINE_SERVICE:
		benchOnlineservice(driver)
	case ARCHIVE_SERVICE:
		benchArchiveService(driver)
	}
}
