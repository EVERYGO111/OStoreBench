package generator

import (
	"github.com/KGXarwen/COSB/distribution"
	"sync"
)

type GeneratorImpl struct {
	distri distribution.Distribution
	sync.Mutex
}

type GeneratorArgs map[string]interface{}

func (args GeneratorArgs) Get(key string, defaultValue interface{}) interface{} {
	if v, ok := args[key]; ok {
		return v
	}
	return defaultValue
}

func NewGeneratorImpl(distr string, args GeneratorArgs) *GeneratorImpl {
	switch distr {
	case "NegativeExp":
		return &GeneratorImpl{distri: distribution.NewNegativeExp(args.Get("lamdda", 0.05).(float64))}
	case "zipf":
		return &GeneratorImpl{distri: distribution.NewZipf(args.Get("s", 1.5).(float64), args.Get("v", 2).(float64), args.Get("imax", 64*1024*1024).(uint64))}
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
