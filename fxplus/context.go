package fxplus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.uber.org/fx"
)

type AsyncRunner interface {
	MustAsync(func() error)
}

type AsyncInit func(func(ctx context.Context) error)

var ErrContextShutdown = errors.New("shutdown")

func IsShutdownOrCancel(err error) bool {
	return errors.Is(err, ErrContextShutdown) || errors.Is(err, context.Canceled)
}

func Context(
	lc fx.Lifecycle,
	s fx.Shutdowner,
	log *slog.Logger,
) (context.Context, AsyncInit) {
	if log == nil {
		log = slog.Default()
	}
	ctx, cn := context.WithCancelCause(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cn(fmt.Errorf("%w: %w", context.Canceled, ErrContextShutdown))
			return nil
		},
	})
	return ctx, func(fn func(ctx context.Context) error) {
		go func() {
			err := fn(ctx)
			if err != nil {
				log.Error("Failed to run async hook", "err", err)
				s.Shutdown()
			}
		}()

	}
}
