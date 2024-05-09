package dbstruct

import "github.com/georgysavva/scany/v2/pgxscan"

type Querier = pgxscan.Querier


type ExecQuerier interface {
	Execer
	Querier
}
