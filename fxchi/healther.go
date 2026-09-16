package fxchi

import (
	"context"
	"net/http"
	"time"

	"github.com/gfx-labs/utilgo/fxplus"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type HealtherHandlerParams struct {
	fx.In
	Healthers []fxplus.Healther `group:"fxplus"`
}

func HealtherHandler(p HealtherHandlerParams) RouteResults {
	return RouteResults{
		Route: func(r chi.Router) {
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				hrs := []*fxplus.HealthReport{}
				for _, v := range p.Healthers {
					if v == nil {
						continue
					}
					ctx, cn := context.WithTimeout(r.Context(), 5*time.Second)
					hrs = append(hrs, fxplus.HealthCheck(ctx, v))
					cn()
				}
				fxplus.RespondHealth(w, hrs...)
			})
		},
	}
}
