package fxriver

import (
	"context"
	"fmt"
	"gfx.cafe/util/go/gotel"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

type TraceWorker[T river.JobArgs] struct {
	worker river.Worker[T]
}

func NewTraceWorker[T river.JobArgs](worker river.Worker[T]) river.Worker[T] {
	o := TraceWorker[T]{
		worker: worker,
	}

	return o
}

func (o TraceWorker[T]) NextRetry(job *river.Job[T]) time.Time {
	return o.worker.NextRetry(job)
}

func (o TraceWorker[T]) Timeout(job *river.Job[T]) time.Duration {
	return o.worker.Timeout(job)
}

func (o TraceWorker[T]) Work(ctx context.Context, job *river.Job[T]) (err error) {
	taskName := job.Args.Kind()
	pkg, _, _ := strings.Cut(taskName, ".")

	tracer := otel.Tracer(pkg)

	ctx, span := tracer.Start(ctx, taskName)
	defer span.End()

	span.SetAttributes(gotel.MakeKeyValue("p.args", job))

	err = o.worker.Work(ctx, job)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(
			attribute.Bool("isError", true),
			attribute.String("error", err.Error()))
		span.RecordError(err, trace.WithAttributes(
			attribute.String("msg", fmt.Sprintf("error in %s: %s", taskName, err.Error()))))
	} else {
		span.SetStatus(codes.Ok, "completed")
	}

	return
}
