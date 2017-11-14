package distribution

import (
	"math"
	"math/rand"
	"time"
)

// Poisson is Used to generate data which fit poisson distribution
type Poisson struct {
	lambda float64
	rd     *rand.Rand
}

// NewPoisson is a function used to create a struct Poisson
func NewPoisson(lambda float64) *Poisson {
	return &Poisson{lambda: lambda, rd: rand.New(rand.NewSource(int64(time.Now().Second())))}
}

// Int64 is used to generate a int64 number
func (ps *Poisson) Int64() int64 {
	L := math.Exp(-ps.lambda)
	k := int64(0)
	p := float64(1)
	for p > L {
		k++
		p *= ps.rd.Float64()
	}
	return k - 1
}
