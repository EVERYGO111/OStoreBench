package workloads

import (
	"fmt"
	"github.com/KDF5000/COSB/distribution"
	"github.com/KDF5000/COSB/driver"
	"math/rand"
	"time"
)

// PoissonWorkload is
type PoissonWorkload struct {
	avgRequestRate float64 //average requests per second used as parameter lambda of poission distribution
	maxFileSize    uint64  //average file size used as parameter lambda of zipf distribution
	count          int64   //the count needed to write to the storage system
	client         driver.Client
}

// NewPoissonWorkloads create
func NewPoissonWorkloads(avgRequestRate float64, maxFileSize uint64,
	cout int64, client driver.Client) *PoissonWorkload {
	return &PoissonWorkload{avgRequestRate: avgRequestRate, maxFileSize: maxFileSize, count: cout, client: client}
}

// Start is
func (w *PoissonWorkload) Start() {
	// file size is zipf distribution
	var fileDataSize []int
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	zipf := rand.NewZipf(r, 2.7, float64(w.maxFileSize/3), w.maxFileSize) //10485760KB = 10 * 1024 * 1024 = 10G

	for i := int64(0); i != w.count; i++ {
		item := int(zipf.Uint64())
		fileDataSize = append(fileDataSize, item)
	}

	//request rate is poisson distribution
	// use exponential distribution describes time between two requests
	p := distribution.NewNegativeExp(1000 / w.avgRequestRate) // average qps is 50
	var waitTimes []int64
	for i := int64(0); i < w.count-1; i++ {
		item := p.Float64()
		waitTimes = append(waitTimes, int64(item)*1000)
	}

	//begin benchmark
	i := int64(0)
	for {
		resps := fmt.Sprintf("%d", fileDataSize[i])
		w.client.SendRequest([]byte(resps))
		if i == w.count-1 {
			break
		}
		time.Sleep(time.Duration(waitTimes[i]) * time.Microsecond)
		i++
	}
}
