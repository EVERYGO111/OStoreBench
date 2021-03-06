package workloads

import (
	"fmt"
	"github.com/KGXarwen/COSB/common"
	"github.com/KGXarwen/COSB/distribution"
	"github.com/KGXarwen/COSB/driver"
	"github.com/KGXarwen/COSB/generator"
	"os"
	"sync"
	"time"
)

const (
	ZIPF = iota
	LOGNORMAL
	POISSION
	NEGATIVE_EXP
)

type DistributionType int

type WorkloadConfig struct {
	WriteRate    float64 //Write request rate, readRate = 1- writeRate
	TotalProcess int64   //total goroutine

	Driver       driver.Driver
	FileSizeType DistributionType
	IatType      DistributionType
	RequestType  DistributionType

	RequestNum int64
}

type Workload struct {
	ReadProcess  int64
	WriteProcess int64
	TotalProcess int64
	RequestNum   int64

	driver            driver.Driver
	fileSizeGenerator generator.Generator
	iatGenerator      generator.Generator
	requestGenerator  generator.Generator //not used yet

	fids   []string
	l      sync.Mutex
	RwLock sync.RWMutex

	RequestPool *ConcurrentPool

	curRequestNum int64
	slock         sync.Mutex //for status update
}

func loadGeneratorByType(t DistributionType) generator.Generator {
	var args generator.GeneratorArgs = make(map[string]interface{})
	switch t {
	case ZIPF:
		args["s"] = float64(1.5)
		args["v"] = float64(2)
		args["imax"] = uint64(64 * 1024 * 1024) //64M
		return generator.NewGeneratorImpl("zipf", args)
	case LOGNORMAL:
		return generator.NewGeneratorImpl("lognormal", args)
	case POISSION:
		return generator.NewGeneratorImpl("poission", args)
	case NEGATIVE_EXP:
		args["lambda"] = float64(0.05)
		return generator.NewGeneratorImpl("NegativeExp", args)
	default:
		args["s"] = float64(1.5)
		args["v"] = float64(2)
		args["imax"] = uint64(64 * 1024 * 1024)
		return generator.NewGeneratorImpl("zipf", args)
	}
}

func NewWorkload(config WorkloadConfig) *Workload {
	w := &Workload{driver: config.Driver}
	w.SetFileSizeGenerator(loadGeneratorByType(config.FileSizeType))
	w.SetIatGenerator(loadGeneratorByType(config.IatType))
	w.SetRequestRateGenerator(loadGeneratorByType(config.RequestType))
	w.TotalProcess = config.TotalProcess
	w.WriteProcess = int64(config.WriteRate * float64(config.TotalProcess))
	w.ReadProcess = config.TotalProcess - w.WriteProcess
	w.RequestNum = config.RequestNum
	w.RequestPool = NewConcurrentPool()

	return w
}

func (w *Workload) SetFileSizeGenerator(g generator.Generator) {
	w.fileSizeGenerator = g
}

func (w *Workload) SetIatGenerator(g generator.Generator) {
	w.iatGenerator = g
}

func (w *Workload) SetRequestRateGenerator(g generator.Generator) {
	w.requestGenerator = g
}

func (w *Workload) appendFid(fid string) {
	w.RwLock.Lock()
	defer w.RwLock.Unlock()
	w.fids = append(w.fids, fid)
}

func (w *Workload) nextFid(zipf distribution.Distribution) string {
	for {
		index := zipf.Uint64()
		//w.RwLock.RLock()
		//fmt.Printf("index:%d, len:%d\n", index, len(w.fids))
		if index >= uint64(len(w.fids)) {
			time.Sleep(1 * time.Microsecond) //if not sleep, there will be a fatal error!!!
			continue
		}
		fid := w.fids[index]
		//w.RwLock.RUnlock()
		return fid
	}
}

func (w *Workload) readFiles(wg *sync.WaitGroup, requestChan chan interface{}, requestDone chan bool, s *stats, localStartIndex int) {
	zipf := distribution.NewZipf(1.5, 2, uint64(w.RequestNum))
	for i := int64(0); i < w.ReadProcess; i++ {
		wg.Add(1)
		go func(st *stat) {
			for {
				select {
				case <-requestChan:
					key := w.nextFid(zipf)
					start := time.Now()
					if bytesRead, err := w.driver.Get(BucketName, key); err == nil {
						st.completed++
						st.transferred += int64(len(bytesRead))
						s.ReadFileIds.Append(key)
						s.ReadFileSize.Append(int64(len(bytesRead)))
						s.ReadTime.Append(time.Now().Sub(start))
					} else {
						st.failed++
						fmt.Printf("Get:%v\n", err)
					}
				case <-requestDone:
					wg.Done()
					return
				}
			}
		}(&s.localStats[localStartIndex+int(i)])
	}
}

func (w *Workload) writeFiles(wg *sync.WaitGroup, requestChan chan interface{}, requestDone chan bool, s *stats, localStartIndex int) {
	for i := int64(0); i < w.WriteProcess; i++ {
		wg.Add(1)
		go func(st *stat) {
			for {
				select {
				case <-requestChan:
					fileName := "cfsb_test_" + time.Now().String()
					fileSize := w.fileSizeGenerator.Uint64() + 1024 // filesize can not be 0
					start := time.Now()
					if fileKey, err := w.driver.Put(BucketName, fileName, int64(fileSize)); err == nil {
						w.appendFid(fileKey)
						st.completed++
						st.transferred += int64(fileSize)
						s.WriteFileSize.Append(fileSize)
						s.WriteTime.Append(time.Now().Sub(start))
					} else {
						st.failed++
						fmt.Printf("Put:%v\n", err)
					}
				case <-requestDone:
					wg.Done()
					return
				}
			}
		}(&s.localStats[localStartIndex+int(i)])
	}
}

func (w *Workload) generateRequest(wg *sync.WaitGroup, requestChan chan interface{}, finishDone chan bool) {
	for i := int64(0); i < w.RequestNum; i++ {
		requestChan <- 1
		t := w.iatGenerator.Uint64()
		time.Sleep(time.Duration(t) * time.Microsecond)
	}
	//totalprocess+1, one for checkprogress
	for i := int64(0); i < w.TotalProcess+1; i++ {
		finishDone <- true
	}
	wg.Done()
}

func (w *Workload) Start() {
	stats := NewStas(int(w.TotalProcess))
	stats.total = w.RequestNum
	stats.start = time.Now()
	//1. iat间隔
	//2. 文件大小
	//3. 读写比例
	wg := &sync.WaitGroup{}
	requestChan := make(chan interface{})
	requestDone := make(chan bool)
	wg.Add(1)
	go w.generateRequest(wg, requestChan, requestDone)
	w.readFiles(wg, requestChan, requestDone, stats, 0)
	w.writeFiles(wg, requestChan, requestDone, stats, int(w.ReadProcess))

	wg.Add(1)
	go stats.CheckProgress("Basic Workload", requestDone, wg)

	wg.Wait()
	stats.end = time.Now()
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
				} else if s, ok := d.Value.(uint64); ok {
					f.WriteString(fmt.Sprintf("%d\n", s))
				}
			}
		}
	})
}
