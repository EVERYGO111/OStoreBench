package distribution

import (
	"math/rand"
	"time"
)

type ZipfDistribution struct {
	zipf *rand.Zipf
}

func NewZipf(s, v float64, imax uint64) *ZipfDistribution {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	return &ZipfDistribution{
		rand.NewZipf(r, s, v, imax),
	}
}

func (z *ZipfDistribution) Uint64() uint64 {
	return z.zipf.Uint64()
}

func (z *ZipfDistribution) Float64() float64 {
	return float64(z.zipf.Uint64())
}

func (z *ZipfDistribution) Int64() int64 {
	return int64(z.zipf.Uint64())
}
