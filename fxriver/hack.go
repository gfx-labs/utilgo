package fxriver

import (
	"github.com/riverqueue/river"
)

type WorkConfigurer interface {
	Configure(workers *river.Workers, useTracing bool)
}

func Wrap[T river.JobArgs](w river.Worker[T]) WorkConfigurer {
	return &wc[T]{worker: w}
}

type wc[T river.JobArgs] struct {
	worker river.Worker[T]
}

func (c *wc[T]) Configure(workers *river.Workers, useTracing bool) {
	if useTracing {
		c.worker = NewTraceWorker[T](c.worker)
	}
	river.AddWorker(workers, c.worker)
}
