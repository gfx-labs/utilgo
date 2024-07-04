package fxriver

import (
	"context"
	"fmt"
	"os"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

type errorHandler struct {
}

// HandleError is invoked in case of an error occurring in a job.
// Context is descended from the one used to start the River client that
// worked the job.
func (errorhandler *errorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
	return &river.ErrorHandlerResult{
		SetCancelled: false,
	}
}

// HandlePanic is invoked in case of a panic occurring in a job.
//
// Context is descended from the one used to start the River client that
// worked the job.
func (errorhandler *errorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any, trace string) *river.ErrorHandlerResult {
	fmt.Fprintf(os.Stderr, "Panic: %v\n%s\n", panicVal, trace)
	return &river.ErrorHandlerResult{
		SetCancelled: false,
	}
}
