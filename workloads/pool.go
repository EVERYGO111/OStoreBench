package workloads

import (
	"container/list"
	"sync"
)

type ConcurrentPool struct {
	sync.Mutex
	pool *list.List
}

func NewConcurrentPool() *ConcurrentPool {
	return &ConcurrentPool{
		pool: list.New(),
	}
}

func (p *ConcurrentPool) Push(v interface{}) {
	if p.pool == nil {
		p.pool = list.New()
	}
	p.Lock()
	defer p.Unlock()
	p.pool.PushBack(v)
}

func (p *ConcurrentPool) Pop() interface{} {
	if p.pool == nil || p.pool.Len() == 0 {
		return nil
	}
	ele := p.pool.Back()
	p.pool.Remove(ele)
	return ele.Value
}
