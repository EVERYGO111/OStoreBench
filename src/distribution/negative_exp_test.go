package distribution

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	exp := NewNegativeExp(0.05)
	for i := 0; i < 10000; i++ {
		fmt.Println(uint64(exp.Float64()))
	}
}
