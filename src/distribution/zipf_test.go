package distribution

import (
	"testing"
	"fmt"
)

func TestZipfBasic(t *testing.T) {
	zipf := NewZipf(1.5,2,64*1024*1024)
	for i := 0; i < 10000; i++ {
		fmt.Println(uint64(zipf.Float64()))
	}
}

