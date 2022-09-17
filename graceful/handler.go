package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// StartFunc is a function that uses a context and returns an error if any bit fails to run. it will be run on a new goroutine using context.background()
// it should not block.
type StartFunc func(ctx context.Context) error

// this is the function that will be run on shutdown. the context will be canclled in DefaultShutdownTimeout after ctrl-c is called
type ShutdownFunc func(context.Context) error

// this is in how long the context will be cancelled.
// you can modify it for your application
var DefaultShutdownTimeout = 5 * time.Second

// handler runs a function that uses a context and then shuts it down  on ctrl-c
// the startfunc should not block
// the shutdown func is run on ctrl-c
// the Handler itself will block
func Handler(start StartFunc, shutdown ShutdownFunc) error {
	var (
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)
	ctx, cancel := context.WithCancel(context.Background())
	// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
	go func() {
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-stopChan
		timer, cn := context.WithTimeout(ctx, DefaultShutdownTimeout)
		defer cancel()
		defer cn()
		if err := shutdown(timer); err != nil {
			errChan <- errors.WithStack(err)
			return
		}

		errChan <- nil
	}()
	if err := start(ctx); err != nil {
		return errors.WithStack(err)
	}
	return <-errChan
}
