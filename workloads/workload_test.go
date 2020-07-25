package workloads

import (
	"fmt"
	"github.com/KGXarwen/COSB/distribution"
	"sync"
	"testing"
)

func TestZipf(t *testing.T) {
	var fileDataSize []int

	zipf := distribution.NewZipf(1.5, 2, 102400)
	for i := int64(0); i < 100000; i++ {
		item := int(zipf.Uint64())
		fileDataSize = append(fileDataSize, item)
	}

	for i := 0; i < len(fileDataSize); i++ {
		//fmt.Println(fileDataSize[i])
	}
}

//验证边写边读，读是否满足zipf分布
func TestRWZipf(t *testing.T) {
	//spawn 10 goroutines, three of them are write operation
	//while the others are read operation
	var fids []int = make([]int, 1000)
	var lock sync.Mutex

	var fid int = 0
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			for {
				lock.Lock()
				fid++
				if fid >= 100000 {
					lock.Unlock()
					//fmt.Println("Write done!")
					wg.Done()
					return
				}
				fids = append(fids, fid)
				lock.Unlock()
			}
		}()
	}

	var readCount int = 0
	var rlock sync.Mutex
	zipf := distribution.NewZipf(1.5, 2, 100000)
	for i := 1; i <= 7; i++ {
		wg.Add(1)
		go func() {
			for {
				index := int(zipf.Uint64())
				if index >= len(fids) {
					continue
				}
				fmt.Println(index)
				rlock.Lock()
				readCount++
				if readCount > 100000 {
					rlock.Unlock()
					wg.Done()
					//fmt.Println("Read Done")
					return
				}
				rlock.Unlock()
			}
		}()
	}
	wg.Wait()
}
