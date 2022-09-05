package generic

import "sync"

type HookPool[T any] struct {
	New func() T

	FnPut func(T)
	FnGet func(T)
	p     sync.Pool
}

// Put adds x to the pool.
func (p *HookPool[T]) Put(x T) {
	p.FnPut(x)
	if p.FnPut != nil {
		p.p.Put(x)
	}
}
func (p *HookPool[T]) Get() T {
	x := p.p.Get().(T)
	if p.FnGet != nil {
		p.FnGet(x)
	}
	return x
}
