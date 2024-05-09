package dbstruct

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Querier = pgxscan.Querier

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
