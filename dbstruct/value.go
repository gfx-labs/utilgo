package dbstruct

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func QueueInsert(batch *pgx.Batch, table string, item any, amend func(string) string) (*pgx.QueuedQuery, error) {
	query, err := MakeQuery(table, item)
	if err != nil {
		return nil, err
	}
	if amend != nil {
		query = amend(query)
	}
	values, err := GetArgs(item)
	if err != nil {
		return nil, err
	}
	return batch.Queue(query, values...), nil
}
func Insert(ctx context.Context, execer Execer, table string, item any, amend func(string) string) (r pgconn.CommandTag, err error) {
	// get
	query, err := MakeQuery(table, item)
	if err != nil {
		return r, err
	}
	if amend != nil {
		query = amend(query)
	}
	values, err := GetArgs(item)
	if err != nil {
		return r, err
	}
	return execer.Exec(ctx, query, values...)
}
