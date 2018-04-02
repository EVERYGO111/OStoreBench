package workloads

import (
	"common"
	"fmt"
	"math"
	"sync"
	"time"
)

type stats struct {
	ReadTime  *common.ConcurrentSlice
	WriteTime *common.ConcurrentSlice

	ReadFileSize  *common.ConcurrentSlice
	ReadFileIds   *common.ConcurrentSlice
	WriteFileSize *common.ConcurrentSlice

	localStats []stat
	start      time.Time
	end        time.Time
	total      int64
}

type stat struct {
	completed   int
	failed      int
	transferred int64
}

func NewStas(n int) *stats {
	return &stats{
		ReadTime:      common.NewConcurrentSlice(),
		WriteTime:     common.NewConcurrentSlice(),
		ReadFileSize:  common.NewConcurrentSlice(),
		ReadFileIds:   common.NewConcurrentSlice(),
		WriteFileSize: common.NewConcurrentSlice(),
		localStats:    make([]stat, n),
	}
}

func (s *stats) CheckProgress(testName string, finishChan chan bool, wg *sync.WaitGroup) {
	fmt.Printf("\n------------ %s ----------\n", testName)
	ticker := time.Tick(time.Second)
	lastCompleted, lastTransferred, lastTime := 0, int64(0), time.Now()
	for {
		select {
		case <-finishChan:
			wg.Done()
			return
		case t := <-ticker:
			completed, transferred, taken, total := 0, int64(0), t.Sub(lastTime), s.total
			for _, localStat := range s.localStats {
				completed += localStat.completed
				transferred += localStat.transferred
			}
			fmt.Printf("Completed %d of %d requests, %3.1f%% %3.1f/s %3.1fMB/s\n",
				completed, total, float64(completed)*100/float64(total),
				float64(completed-lastCompleted)*float64(int64(time.Second))/float64(int64(taken)),
				float64(transferred-lastTransferred)*float64(int64(time.Second))/float64(int64(taken))/float64(1024*1024),
			)
			lastCompleted, lastTransferred, lastTime = completed, transferred, t
		}
	}
}

func (s *stats) PrintStats(process func(times *common.ConcurrentSlice, tag string)) {
	completed, failed, transferred := 0, 0, int64(0)
	for _, ls := range s.localStats {
		completed += ls.completed
		failed += ls.failed
		transferred += ls.transferred
	}

	timeToke := float64(s.end.Sub(s.start)) / 1000000000
	fmt.Printf("\nTotal Request Process:    %d\n", len(s.localStats))
	fmt.Printf("Time token for test:        %.3f seconds\n", timeToke)
	fmt.Printf("Complete Requests:          %d\n", completed)
	fmt.Printf("Failed Requests:            %d\n", failed)
	fmt.Printf("Total transferred:          %d bytes\n", transferred)
	fmt.Printf("Request per second:         %.2f [#/sec]\n", float64(completed)/timeToke)
	fmt.Printf("Transfer rate:              %.2f [Kbytes/sec]\n", float64(transferred)/1024/timeToke)

	if process != nil {
		process(s.ReadTime, "read_time")
		process(s.WriteTime, "write_time")
	}
	if s.ReadTime.Len() > 0 {
		printTime(s.ReadTime, "Read Times(ms)")
	}

	if s.WriteTime.Len() > 0 {
		printTime(s.WriteTime, "Write Times(ms)")
	}
}

func printTime(slice *common.ConcurrentSlice, tips string) {
	if slice == nil {
		return
	}
	sum, min, max := float64(0), float64(math.MaxInt64), float64(0)
	for d := range slice.Iter() {
		if t, ok := d.Value.(time.Duration); ok {
			dMs := float64(t) / 1000000
			if dMs < min {
				min = dMs
			}
			if dMs > max {
				max = dMs
			}
			sum += dMs
		} else {
			fmt.Printf("Readtime format error!")
			continue
		}
	}
	avg := sum / float64(slice.Len())
	varianceSum := 0.0
	for d := range slice.Iter() {
		if t, ok := d.Value.(time.Duration); ok {
			dMs := float64(t) / 1000000
			diff := float64(dMs) - avg
			varianceSum += diff * diff

		} else {
			fmt.Printf("Readtime format error!")
			continue
		}
	}
	std := math.Sqrt(varianceSum / float64(slice.Len()))
	fmt.Printf("\n%s\n", tips)
	fmt.Printf("           min       avg        max       std\n")
	fmt.Printf("Total:     %.2f      %.2f       %.2f      %.2f\n", min, avg, max, std)
}
