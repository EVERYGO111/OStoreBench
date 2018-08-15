package distribution

type Distribution interface {
	Uint64() uint64
	Float64() float64
	Int64() int64
}
