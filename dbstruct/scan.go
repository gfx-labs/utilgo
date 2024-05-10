package dbstruct

import (
	"context"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

var scany_api, _ = dbscan.NewAPI(dbscan.WithAllowUnknownColumns(true))
var Scan, _ = pgxscan.NewAPI(scany_api)

func Get(ctx context.Context, db Querier, dst any, query string, args ...any) error {
	return Scan.Get(ctx, db, dst, query, args...)
}
func Select(ctx context.Context, db Querier, dst any, query string, args ...any) error {
	return Scan.Select(ctx, db, dst, query, args...)
}
func ScanAll(dst any, rows pgx.Rows) error {
	return Scan.ScanAll(dst, rows)
}

func ScanOne(dst any, rows pgx.Rows) error {
	return Scan.ScanOne(dst, rows)
}
