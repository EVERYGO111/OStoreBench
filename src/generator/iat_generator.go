package generator

import "distribution"

type IATGenerator struct {
	distri distribution.Distribution
}

func NewIATGenerator(distr string) *IATGenerator {
	switch distr {
	case "NegativeExp":
		return &IATGenerator{distri: distribution.NewNegativeExp(50)}
	}
	return nil
}

func (g *IATGenerator) Uint64() uint64 {
	if g.distri != nil {
		return g.distri.Uint64()
	}
	return 0
}
