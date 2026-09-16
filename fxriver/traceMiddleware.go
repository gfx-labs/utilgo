package fxriver

import (
	"context"
	"fmt"
	"strings"

	"github.com/gfx-labs/utilgo/gotel"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewTraceMiddleware() rivertype.WorkerMiddleware {
	return &traceWorkerMiddleware{}
}

type traceWorkerMiddleware struct {
	// embed JobInsertMiddlewareDefaults for forward compatibility
	// in case additional methods are added to the interface:
	river.WorkerMiddlewareDefaults
}

func (m *traceWorkerMiddleware) Work(ctx context.Context, job *rivertype.JobRow, doInner func(ctx context.Context) error) error {
	// Extract the trace ID from the job metadata and log it.
	taskName := job.Kind
	pkg, _, _ := strings.Cut(taskName, ".")

	tracer := otel.Tracer(pkg)

	ctx, span := tracer.Start(ctx, taskName)
	defer span.End()

	span.SetAttributes(gotel.MakeKeyValue("p.args", job))

	err := doInner(ctx)
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
	return err
}
