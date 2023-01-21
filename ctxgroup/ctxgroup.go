package ctxgroup

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// a Group is an errgroup that also has a context in it
// this ctx will be used as the context of the waitgroup
type Group struct {
	ctx context.Context
	g   *errgroup.Group

	preinit bool
	o       sync.Once
}

func (g *Group) init() {
	if !g.preinit {
		g.o.Do(func() {
			if g.g == nil {
				g.g, g.ctx = errgroup.WithContext(context.TODO())
			}
		})
	}
}

func New(ctx context.Context) *Group {
	g, ctx := errgroup.WithContext(ctx)
	return &Group{
		preinit: true,
		g:       g,
		ctx:     ctx,
	}
}

func (g *Group) SetLimit(n int) {
	g.init()
	g.g.SetLimit(n)
}

func (g *Group) Wait() error {
	g.init()
	return g.g.Wait()
}

func (g *Group) Go(fn func(context.Context) error) {
	g.init()
	g.g.Go(func() error {
		return fn(g.ctx)
	})
}

func (g *Group) TryGo(fn func(context.Context) error) bool {
	g.init()
	return g.g.TryGo(func() error {
		return fn(g.ctx)
	})
}
