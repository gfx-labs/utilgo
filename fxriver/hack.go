package fxriver

import (
	"github.com/riverqueue/river"
)

type WorkConfigurer interface {
	Configure(workers *river.Workers)
}

func _Wrap[T river.JobArgs](w river.Worker[T]) WorkConfigurer {
	return &wc[T]{worker: w}
}

func Wrap[T river.JobArgs](w river.Worker[T]) WorkConfigurer {
	return _Wrap(NewTraceWorker[T](w))
}

type wc[T river.JobArgs] struct {
	worker river.Worker[T]
}

func (c *wc[T]) Configure(workers *river.Workers) {
	river.AddWorker(workers, c.worker)
}
