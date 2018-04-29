package workloads

import (
	"common"
	"container/list"
	"driver"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

//
//import (
//	"sync"
//	"time"
//	"fmt"
//	"os"
//	"driver"
//)
//
//// for bigdata processing type workload
//// in the bigdata process scenario, read , write and delete operations
//// always thurn. And the different type operation always appears in
//// different stage.
//
//
type MRWorkloadConfig struct {
	Stages       int //read/write stages
	MaxFileSize  int64
	MaxOpInStage int
	ProcessNum   int
}

type MRWorkload struct {
	driver driver.Driver
	config *MRWorkloadConfig

	deleteQueue *list.List
}

//
func NewMRWorkload(driver driver.Driver, conf *MRWorkloadConfig) *MRWorkload {
	if conf.Stages == 0 {
		conf.Stages = 50
	}
	w := &MRWorkload{config: conf}
	w.driver = driver
	w.deleteQueue = list.New()
	return w
}

func (w *MRWorkload) deleteFiles(files []string, wg *sync.WaitGroup, s *stats, index int) {
	if files == nil {
		wg.Done()
		return
	}
	for _, file := range files {
		start := time.Now()
		err := w.driver.Delete(BucketName, file)
		if err != nil {
			s.localStats[index].failed++
			continue
		}
		s.DeleteTime.Append(time.Now().Sub(start))
	}
	wg.Done()
}

//
func (w *MRWorkload) deleteProcess(deleteWg *sync.WaitGroup, finish chan bool, deleteSignal chan interface{}, s *stats, index int) {
	wg := &sync.WaitGroup{}
	for {
		select {
		case <-finish:
			wg.Wait()
			deleteWg.Done()
			return
		case <-deleteSignal:
			if w.deleteQueue.Len() > 0 {
				wg.Add(1)
				ele := w.deleteQueue.Front()
				files := ele.Value.([]string)
				w.deleteQueue.Remove(ele)
				go w.deleteFiles(files, wg, s, index)
			}
		}
	}

}

func (w *MRWorkload) doDAGRequest(wg *sync.WaitGroup, s *stats, index int) {
	deleteSignal := make(chan interface{})
	finish := make(chan bool)

	deleteWg := &sync.WaitGroup{}
	deleteWg.Add(1)
	go w.deleteProcess(deleteWg, finish, deleteSignal, s, index)

	opType := uint8(0) //0: wirte, 1: read
	writeQueue := make([]string, 0)
	for i := 0; i < w.config.Stages; i++ {
		r := rand.New(rand.NewSource(int64(time.Now().Second())))
		opNum := r.Intn(w.config.MaxOpInStage) + 10
		switch opType {
		case 0:
			//write
			deleteSignal <- 1 //delete old files
			for j := 0; j < opNum; j++ {
				fileSize := r.Int63n(w.config.MaxFileSize)
				fileName := fmt.Sprintf("cfsb_archive_%d.txt", time.Now().UnixNano())
				start := time.Now()
				fileKey, err := w.driver.Put(BucketName, fileName, fileSize)
				if err != nil {
					s.localStats[i].failed++
					continue
				}
				writeQueue = append(writeQueue, fileKey)
				s.WriteTime.Append(time.Now().Sub(start))
				s.localStats[index].transferred += fileSize
				s.WriteFileSize.Append(fileSize)
			}
		case 1:
			//read
			if len(writeQueue) <= 0 {
				continue
			}
			for j := 0; j < opNum; j++ {
				fileKey := writeQueue[r.Intn(len(writeQueue))]
				start := time.Now()
				bytesRead, err := w.driver.Get(BucketName, fileKey)
				if err != nil {
					s.localStats[index].failed++
					continue
				}
				s.ReadTime.Append(time.Now().Sub(start))
				s.localStats[index].transferred += int64(len(bytesRead))
			}
			oldWriteQueue := writeQueue
			w.deleteQueue.PushBack(oldWriteQueue)
			writeQueue = make([]string, 0)
		}
		opType = opType ^ 0x01
		s.localStats[index].completed++
	}
	finish <- true
	deleteWg.Wait()
	wg.Done()
}

//
// spawn w.config.ProcessNum goroutine,
// each process will send w.config.requestGroups group requests,
// which contains random count reqeuests in one group
func (w *MRWorkload) Start() {
	stats := NewStas(int(w.config.ProcessNum))
	stats.total = int64(w.config.ProcessNum * w.config.Stages)
	stats.start = time.Now()

	wg := &sync.WaitGroup{}
	for i := 0; i < w.config.ProcessNum; i++ {
		wg.Add(1)
		go w.doDAGRequest(wg, stats, i)
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
					f.WriteString(fmt.Sprintf("%d", t))
				}
			}
		}
	})
}
