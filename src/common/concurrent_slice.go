package common

import "sync"

type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

func NewConcurrentSlice() *ConcurrentSlice {
	return &ConcurrentSlice{
		items: make([]interface{}, 0),
	}
}

func (s *ConcurrentSlice) Append(value interface{}) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, value)
}

func (s *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)
	go func() {
		s.RLock()
		defer s.RUnlock()
		for index, v := range s.items {
			c <- ConcurrentSliceItem{Index: index, Value: v}
		}
		close(c)
	}()
	return c
}


func (s *ConcurrentSlice)Len()int{
	return len(s.items)
}