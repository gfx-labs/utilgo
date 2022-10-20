package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	Cancelled         = cancelReason("context cancelled")
	Stopped           = cancelReason("retry stopped")
	AttemptsExhausted = cancelReason("attempts exhausted")
)

type Func func(context.Context, int) error

type cancelReason string

type stopErr struct {
	err error
}

func (e *stopErr) Error() string {
	return fmt.Sprintf("stop err: %s", e.err.Error())
}

type Err struct {
	Err      error
	Reason   cancelReason
	Attempts int
}

func (e *Err) Cause() error { return e.Err }
func (e *Err) Error() string {
	return fmt.Sprintf("on attempt '%d'; %s: %s", e.Attempts, e.Reason, e.Err.Error())
}

func (e *Err) Is(target error) bool {
	_, ok := target.(*Err)
	return ok
}

// Stop forces the retry to cancel with the provided error
// and retry.Err.Reason == retry.Stopped
func Stop(err error) error {
	return &stopErr{err: err}
}

// Until will retry the provided `retry.Func` until it returns nil or
// the context is cancelled. Optionally users may use `retry.Stop()` to force
// the retry to terminate with an error. Returns a `retry.Err` with
// the included Reason and Attempts
func Until(ctx context.Context, backOff BackOff, f Func) error {
	var attempt int
	for {
		attempt++
		if err := f(ctx, attempt); err != nil {
			var stop *stopErr
			if errors.As(err, &stop) {
				return &Err{Attempts: attempt, Reason: Stopped, Err: stop.err}
			}
			interval, retry := backOff.Next()
			if !retry {
				return &Err{Attempts: attempt, Reason: AttemptsExhausted, Err: err}
			}
			timer := time.NewTimer(interval)
			select {
			case <-timer.C:
				timer.Stop()
				continue
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return &Err{Attempts: attempt, Reason: Cancelled, Err: err}
			}
		}
		return nil
	}
}
