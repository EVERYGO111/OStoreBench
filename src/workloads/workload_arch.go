// for archive type workload
// in the archive scenario, there are always be dense write operations periodically
//
package workloads

import (
	"driver"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ArchiveWorkloadConfig struct {
	RequestGroups    int   //total groups
	GroupInteralTime int   //second
	MinFileSize      int64 //bytes
	MaxFileSize      int64 //bytes
	MaxFilesPerGroup int

	ProcessNum int
}

type ArchWorkload struct {
	driver driver.Driver
	config *ArchiveWorkloadConfig
}

func NewArchWorkload(driver driver.Driver, conf *ArchiveWorkloadConfig) *ArchWorkload {
	if conf.ProcessNum == 0 {
		conf.ProcessNum = 16
	}
	w := &ArchWorkload{config: conf}
	w.driver = driver
	return w
}

func (w *ArchWorkload) generateGroups(requestGroup chan interface{}) {
	for i := 0; i < w.config.RequestGroups; i++ {
		requestGroup <- 1
	}
}

func (w *ArchWorkload) sendRequest(requestChan chan interface{}, wg *sync.WaitGroup) {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	rqNum := r.Intn(w.config.MaxFilesPerGroup) + 1
	for i := 0; i < rqNum; i++ {
		fileSize := r.Int63n(w.config.MaxFileSize-w.config.MinFileSize) + w.config.MinFileSize
		fileName := fmt.Sprintf("weed_archive_%d.txt", time.Now().UnixNano())
		fileKey, err := w.driver.Put(fileName, fileSize)
		if err != nil {
			fmt.Printf("Write file error:%v\n", err)
			continue
		}
		fmt.Printf("Put file %v\n", fileKey)
	}
	wg.Done()
}

func (w *ArchWorkload) Start() {
	var wg sync.WaitGroup
	requestGroupChan := make(chan interface{})

	for i := 0; i < w.config.ProcessNum; i++ {
		wg.Add(1)
		go w.sendRequest(requestGroupChan, &wg)
		time.Sleep(time.Duration(w.config.GroupInteralTime) * time.Second)
	}

	wg.Wait()
}
