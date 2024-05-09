package dbstruct

import (
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var scany_api, _ = dbscan.NewAPI(dbscan.WithAllowUnknownColumns(true))
var Scan, _ = pgxscan.NewAPI(scany_api)
