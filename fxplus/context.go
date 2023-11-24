package fxplus

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

type AsyncRunner interface {
	MustAsync(func() error)
}

type AsyncInit func(func(ctx context.Context) error)

func Context(
	lc fx.Lifecycle,
	s fx.Shutdowner,
	log *slog.Logger,
) (context.Context, AsyncInit) {
	if log == nil {
		log = slog.Default()
	}
	ctx, cn := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cn()
			return nil
		},
	})
	return ctx, func(fn func(ctx context.Context) error) {
		go func() {
			err := fn(ctx)
			if err != nil {
				log.Error("Failed to run async hook", err)
				s.Shutdown()
			}
		}()

	}
}
