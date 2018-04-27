package driver

import (
	"fmt"
	"testing"
	"time"
)

func TestCephBasic(t *testing.T) {
	driver := NewCephDriver(
		"testuser:swift",
		"KZMKDDmgJl9dCSkFQi1AzBAYLjNDSOwMqaqfCxM1", //password
		"http://172.16.1.92:5000/v3",
	)

	fileName := fmt.Sprintf("ceph_%d.txt", time.Now().UnixNano())
	fileKey, err := driver.Put("test", fileName, 1024)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put file %s\n", fileKey)
	//get file
	data, err := driver.Get("test", fileKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get file %s ,length:%d\n", fileKey, len(data))
}
