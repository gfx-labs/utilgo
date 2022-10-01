package generic

import "sync"

type HookPool[T any] struct {
	New func() T

	FnPut func(T)
	FnGet func(T)
	p     sync.Pool
	o     sync.Once
}

// Put adds x to the pool.
func (p *HookPool[T]) Put(x T) {
	if p.FnPut != nil {
		p.FnPut(x)
	}
	p.p.Put(x)
}
func (p *HookPool[T]) Get() T {
	p.o.Do(func() {
		p.p.New = func() any {
			return p.New()
		}
	})
	x := p.p.Get().(T)
	if p.FnGet != nil {
		p.FnGet(x)
	}
	return x
}
