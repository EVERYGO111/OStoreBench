// for archive type workload
// in the archive scenario, there are always be dense write operations periodically
//
package workloads

import (
	"github.com/KGXarwen/COSB/common"
	"github.com/KGXarwen/COSB/driver"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	BucketName = "container1"
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

func (w *ArchWorkload) sendRequest(wg *sync.WaitGroup, s *stats, index int) {
	for i := 0; i < w.config.RequestGroups; i++ {
		r := rand.New(rand.NewSource(int64(time.Now().Second())))
		rqNum := r.Intn(w.config.MaxFilesPerGroup) + 1
		for i := 0; i < rqNum; i++ {
			fileSize := r.Int63n(w.config.MaxFileSize-w.config.MinFileSize) + w.config.MinFileSize
			fileName := fmt.Sprintf("cfsb_archive_%d.txt", time.Now().UnixNano())
			start := time.Now()
			_, err := w.driver.Put(BucketName, fileName, fileSize)
			if err != nil {
				fmt.Printf("Write file error:%v\n", err)
				s.localStats[i].failed++
				continue
			}
			s.WriteTime.Append(time.Now().Sub(start))
			s.localStats[index].transferred += fileSize
			s.WriteFileSize.Append(fileSize)
			//fmt.Printf("Put file %v\n", fileKey)
		}
		s.localStats[index].completed++
	}
	wg.Done()
}

// spawn w.config.ProcessNum goroutine,
// each process will send w.config.requestGroups group requests,
// which contains random count reqeuests in one group
func (w *ArchWorkload) Start() {
	stats := NewStas(int(w.config.ProcessNum))
	stats.total = int64(w.config.RequestGroups * w.config.ProcessNum)
	stats.start = time.Now()

	var wg sync.WaitGroup
	for i := 0; i < w.config.ProcessNum; i++ {
		wg.Add(1)
		go w.sendRequest(&wg, stats, i)
	}

	finishDone := make(chan bool)
	svg := &sync.WaitGroup{}
	svg.Add(1)
	go stats.CheckProgress("Arch workload", finishDone, svg)
	wg.Wait()
	stats.end = time.Now()

	finishDone <- true
	svg.Wait()

	stats.PrintStats(func(times *common.ConcurrentSlice, tag string) {
		if times != nil {
			f, err := os.Create(tag)
			defer f.Close()
			if err != nil {
				fmt.Printf("%s: %v\n", err)
				return
			}
			for d := range times.Iter() {
				if t, ok := d.Value.(time.Duration); ok {
					f.WriteString(fmt.Sprintf("%d\n", t))
				}else if s, ok := d.Value.(int64); ok {
					f.WriteString(fmt.Sprintf("%d\n", s))
				}
			}
		}
	})
}
