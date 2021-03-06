package distribution

import (
	"math"
	"math/rand"
	"time"
)

// Exp is used to generate numbers whose distribution is the exponential distribution
type Exp struct {
	lambdar float64
	rd      *rand.Rand
}

// NewNegativeExp create struct Exp
// lambdar =  1 / lambda
func NewNegativeExp(lambda float64) *Exp {
	return &Exp{lambdar: 1/lambda, rd: rand.New(rand.NewSource(int64(time.Now().Second())))}
}

// Float64 return a float64 number
func (e *Exp) Float64() float64 {
	var z = float64(0)
	for z == 0 || z == 1 {
		z = e.rd.Float64()
	}
	//fmt.Printf("z:%f,log:%f\n", 1-z,math.Log(1-z))
	return -e.lambdar * math.Log(1-z)
}

func (e *Exp) Int64() int64 {
	return int64(e.Float64())
}

func (e *Exp) Uint64() uint64 {
	return uint64(e.Float64())
}
