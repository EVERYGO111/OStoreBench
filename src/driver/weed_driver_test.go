package driver

import (
	"testing"
	"fmt"
)

func TestBasic(t *testing.T) {
	driver := NewWeeDriver("localhost:9333")
	fid, err := driver.Put("","test.txt", 1024)
	if err != nil {
		t.Fatalf("Put: %v", err)
	}
	fmt.Printf("Put file, fid:%s\n", fid)
	//t.Logf("Put file, fid:%s\n", fid)
	//get file by fid
	var readBytes []byte
	if readBytes, err = driver.Get("",fid); err != nil {
		t.Fatalf("Get: %v", err)
	}
	fmt.Printf("Get file, fid:%s, length:%d\n", fid, len(readBytes))
	//t.Logf("Get file, fid:%s, length:%d\n", fid, len(readBytes))
	//delete file by id
	if err := driver.Delete("",fid); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	fmt.Printf("Delete file :%s\n", fid)
	//t.Logf("Delete file :%s\n", fid)
}
