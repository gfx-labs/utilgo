package fxplus

import (
	"fmt"
	"net/http"

	"go.uber.org/fx"
)

func ServerWithName(fn func() *http.Server, name string) any {
	return fx.Annotate(fn, fx.ResultTags(fmt.Sprintf(`name:"%s"`, name)))
}
