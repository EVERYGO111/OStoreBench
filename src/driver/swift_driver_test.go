package driver

import (
	"fmt"
	"testing"
	"time"
)

func TestSwiftBasic(t *testing.T) {
	driver := NewSwiftDriver(
		"demo",
		"123456", //password
		"http://172.16.1.92:5000/v3",
		"default",
		"demo", //project name in V3
	)

	fileName := fmt.Sprint("swift_%d.txt", time.Now().UnixNano())
	fileKey, err := driver.Put("container1", fileName, 1024)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put file %s\n", fileKey)
	//get file
	data, err := driver.Get("container1", fileKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get file %s ,length:%d\n", fileKey, len(data))

}
