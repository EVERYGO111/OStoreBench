package generator

import (
	"distribution"
	"sync"
)

type GeneratorImpl struct {
	distri distribution.Distribution
	sync.Mutex
}

func NewGeneratorImpl(distr string) *GeneratorImpl {
	switch distr {
	case "NegativeExp":
		return &GeneratorImpl{distri: distribution.NewNegativeExp(0.05)}
	case "zipf":
		return &GeneratorImpl{distri: distribution.NewZipf(1.5, 2, 64*1024*1024)}
	}
	return nil
}

func (g *GeneratorImpl) Uint64() uint64 {
	g.Lock()
	defer g.Unlock()
	if g.distri != nil {
		return g.distri.Uint64()
	}
	return 0
}
