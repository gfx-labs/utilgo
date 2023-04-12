package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// StartFunc is a function that uses a context and returns an error if any bit fails to run. it will be run on a new goroutine using context.background()
// it should not block.
type StartFunc func(ctx context.Context, done <-chan struct{}) error

// this is the function that will be run on shutdown. the context will be canclled in DefaultShutdownTimeout after ctrl-c is called
type ShutdownFunc func(context.Context) error

// handler runs a function that uses a context and then shuts it down  on ctrl-c
// the startfunc should not block
// the shutdown func is run on ctrl-c
// the Handler itself will block
func Handler(shutdownTime time.Duration, start StartFunc, shutdown ShutdownFunc) error {
	var (
		doneChan = make(chan struct{})
		stopChan = make(chan os.Signal, 1)
		errChan  = make(chan error)
	)
	if start == nil {
		start = func(ctx context.Context, done <-chan struct{}) error {
			return nil
		}
	}
	if shutdown == nil {
		shutdown = func(ctx context.Context) error {
			return nil
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
	go func() {
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-stopChan
		close(doneChan)
		timer, cn := context.WithTimeout(ctx, shutdownTime)
		defer cancel()
		defer cn()
		if err := shutdown(timer); err != nil {
			errChan <- errors.WithStack(err)
			return
		}
		errChan <- nil
	}()
	if err := start(ctx, doneChan); err != nil {
		return errors.WithStack(err)
	}
	return <-errChan
}

func WgCh(s *sync.WaitGroup) chan struct{} {
	o := make(chan struct{})
	go func() {
		s.Wait()
		close(o)
	}()
	return o
}
