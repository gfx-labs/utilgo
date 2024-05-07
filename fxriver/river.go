package fxriver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"go.uber.org/fx"
)

func newRiverConn(ctx context.Context, pgxConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	riverPgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}
	_, err = riverPgxPool.Exec(ctx, `create schema if not exists river0`)
	if err != nil {
		return nil, err
	}
	riverPgx := riverpgxv5.New(riverPgxPool)
	migrator := rivermigrate.New(riverPgx, nil)
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		return nil, err
	}
	return riverPgxPool, nil
}

// scheduler provides a river.client with no name and a *pgxpool.Pool tagged with the name
func Scheduler(name string) any {
	return fx.Annotate(
		SchedulerProvider,
		fx.ParamTags(
			"",
			fmt.Sprintf(`optional:"true"`),
			"",
			//
			fmt.Sprintf(`name:"%s" optional:"true"`, name),
			fmt.Sprintf(`name:"%s"`, name),
		),
		fx.ResultTags(
			``,
			fmt.Sprintf(`name:"%s"`, name),
		),
	)
}

func SchedulerProvider(
	ctx context.Context,
	log *slog.Logger,
	lc fx.Lifecycle,
	//
	config *river.Config,
	pgxConfig *pgxpool.Config,
) (*river.Client[pgx.Tx], *pgxpool.Pool, error) {
	riverPgxPool, err := newRiverConn(ctx, pgxConfig)
	if err != nil {
		return nil, nil, err
	}
	if config == nil {
		config = &river.Config{}
	}
	// default logger is provided slogger
	if config.Logger == nil {
		config.Logger = log
	}
	riverPgx := riverpgxv5.New(riverPgxPool)
	riverClient, err := river.NewClient(riverPgx, config)
	if err != nil {
		return nil, nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return riverClient.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return riverClient.Stop(ctx)
		},
	})
	return riverClient, riverPgxPool, nil
}

// provides a *river.Client[pgx.Tx] with the tag in the name
// consumes WorkConfigurer in a group with the name
// you must provide a *river.Config and *pgxpool.Config with the name
// it will also invoke the client.
func WorkGroupModule(name string, extra ...fx.Option) any {
	return fx.Module("workgroup/"+name,
		append(extra, fx.Invoke(
			fx.Annotate(
				func(*river.Client[pgx.Tx]) {},
				fx.ParamTags(
					fmt.Sprintf(`name:"%s"`, name),
				),
			),
		),
			fx.Provide(
				fx.Annotate(
					WorkGroupProvider,
					fx.ParamTags(
						"",
						fmt.Sprintf(`optional:"true"`),
						"",
						//
						fmt.Sprintf(`name:"%s" optional:"true"`, name),
						fmt.Sprintf(`name:"%s"`, name),
						fmt.Sprintf(`name:"%s"`, name),
						fmt.Sprintf(`group:"%s"`, name),
					),
					fx.ResultTags(
						fmt.Sprintf(`name:"%s"`, name),
					),
				),
			),
		)...)
}

func WorkGroupProvider(
	ctx context.Context,
	log *slog.Logger,
	lc fx.Lifecycle,
	//
	config *river.Config,
	pgxConfig *pgxpool.Config,
	queues map[string]river.QueueConfig,
	workers []WorkConfigurer,
) (*river.Client[pgx.Tx], error) {
	riverPgxPool, err := newRiverConn(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &river.Config{}
	}
	if config.Workers == nil {
		config.Workers = river.NewWorkers()
	}
	// default logger is provided slogger
	if config.Logger == nil {
		config.Logger = log
	}
	for _, v := range workers {
		v.Configure(config.Workers)
	}
	riverPgx := riverpgxv5.New(riverPgxPool)
	riverClient, err := river.NewClient(riverPgx, config)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return riverClient.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return riverClient.Stop(ctx)
		},
	})
	return riverClient, nil
}
