package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func parseTime(date string) int64 {
	timeLayout := "2006/01/02 15:04:05"
	s := strings.Split(date, " ")[:2]
	toBeCharge := strings.Join(s, " ")
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	return t.Unix()
}

func main() {
	fp, err := os.Open("../server/server.log")
	if err != nil {
		log.Fatal(err)
		return
	}
	ioReader := bufio.NewReader(fp)
	beginTime := int64(0)
	currentCount := int(0)
	var data []int
	interval := int64(1)

	for {
		a, _, c := ioReader.ReadLine()
		if c == io.EOF {
			break
		}
		tms := parseTime(string(a))
		// fmt.Println(tms)
		if beginTime == 0 {
			beginTime = tms
			currentCount++
		} else {
			if tms-beginTime < interval {
				currentCount++
			} else {
				data = append(data, currentCount)
				beginTime = tms
				currentCount = 1
			}
		}
	}
	data = append(data, currentCount)
	for i := 0; i < len(data); i++ {
		fmt.Println(data[i])
	}
	defer fp.Close()
}
