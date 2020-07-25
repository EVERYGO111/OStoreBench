package generator

import (
	"github.com/KGXarwen/COSB/distribution"
	"sync"
)

type IATGenerator struct {
	distri distribution.Distribution
	sync.Mutex
}

func NewIATGenerator(distr string) *IATGenerator {
	switch distr {
	case "NegativeExp":
		return &IATGenerator{distri: distribution.NewNegativeExp(50)}
	}
	return nil
}

func (g *IATGenerator) Uint64() uint64 {
	g.Lock()
	defer g.Lock()
	if g.distri != nil {
		return g.distri.Uint64()
	}
	return 0
}
