package dbstruct

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type Querier = pgxscan.Querier

type RowQuerier interface {
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
}

type Batcher interface {
	SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults)
}

type ExecQuerier interface {
	Execer
	Querier
}

type ExecBatchQuerier interface {
	Execer
	Batcher
	Querier
}

type ExecRowQuerier interface {
	ExecQuerier
	RowQuerier
}

type ExecBatchRowQuerier interface {
	ExecBatchQuerier
	RowQuerier
}
