package fxplus

import (
	"fmt"

	"go.uber.org/fx"
)

func WithName(fn any, name string) any {
	return fx.Annotate(fn, fx.ResultTags(fmt.Sprintf(`name:"%s"`, name)))
}
