package workloads

import (
	"distribution"
	"driver"
	"fmt"
	"generator"
	"sync"
	"time"
)

const (
	ZIPF = iota
	LOGNORMAL
	POISSION
	NEGATIVE_EXP
)

var readCnt int
var writeCnt int
var rlock sync.Mutex
var wlock sync.Mutex

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
	switch t {
	case ZIPF:
		return generator.NewGeneratorImpl("zipf")
	case LOGNORMAL:
		return generator.NewGeneratorImpl("lognormal")
	case POISSION:
		return generator.NewGeneratorImpl("poission")
	case NEGATIVE_EXP:
		return generator.NewGeneratorImpl("NegativeExp")
	default:
		return generator.NewGeneratorImpl("zipf")
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

func (w *Workload) readFiles(wg *sync.WaitGroup, requestChan chan interface{}, requestDone chan bool) {
	zipf := distribution.NewZipf(1.5, 2, uint64(w.RequestNum))
	for i := int64(1); i <= w.ReadProcess; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-requestChan:
					if _, err := w.driver.Get(w.nextFid(zipf)); err != nil {
						fmt.Printf("%v\n", err)
					}
					rlock.Lock()
					readCnt++
					rlock.Unlock()
				case <-requestDone:
					wg.Done()
					return
				}
			}
		}()
	}
}

func (w *Workload) writeFiles(wg *sync.WaitGroup, requestChan chan interface{}, requestDone chan bool) {
	for i := int64(1); i <= w.WriteProcess; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-requestChan:
					fileKey, err := w.driver.Put("weed_test_"+time.Now().String(), int64(w.fileSizeGenerator.Uint64()))
					if err != nil {
						fmt.Printf("%v\n", err)
					} else {
						//fmt.Printf("%v\n", fileKey)
						w.appendFid(fileKey)
					}
					wlock.Lock()
					writeCnt++
					wlock.Unlock()
				case <-requestDone:
					wg.Done()
					return
				}
			}
		}()
	}
}

func (w *Workload) generateRequest(wg *sync.WaitGroup, requestChan chan interface{}) {
	for i := int64(0); i < w.RequestNum; i++ {
		requestChan <- 1
		t := w.iatGenerator.Uint64()
		fmt.Println(t)
		time.Sleep(time.Duration(t) * time.Microsecond)
	}
	//totalprocess+1
	for i := int64(0); i <= w.TotalProcess; i++ {
		wg.Done()
	}
}
func (w *Workload) Start() {
	//1. iat间隔
	//2. 文件大小
	//3. 读写比例
	readCnt = 0
	writeCnt = 0
	fmt.Printf("start workload\n")
	wg := &sync.WaitGroup{}
	requestChan := make(chan interface{})
	requestDone := make(chan bool)
	wg.Add(1)
	go w.generateRequest(wg, requestChan)
	w.readFiles(wg, requestChan, requestDone)
	w.writeFiles(wg, requestChan, requestDone)

	wg.Wait()
	fmt.Printf("ReadProcess:%d, WriteProcess:%d, ReadCount:%d, WriteCount:%d\n", w.ReadProcess, w.WriteProcess, readCnt, writeCnt)
}
