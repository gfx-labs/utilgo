package fxmetrics

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func MetricsApp(opts ...fx.Option) fx.Option {
	return fx.Module("fxmetrics",
		fx.Provide(fx.Annotate(http.NewServeMux, fx.ResultTags(`name:"fxmetrics"`))),
		fx.Invoke(Monitoring),
		fx.Invoke(ServerInvoker),
		fx.Options(opts...),
	)
}

type MonitoringParams struct {
	fx.In

	Mux *http.ServeMux `name:"fxmetrics"`
	Lc  fx.Lifecycle
}

func Monitoring(p MonitoringParams) {
	p.Mux.Handle("/debug", middleware.Profiler())
	p.Mux.Handle("/metrics", promhttp.Handler())
}

type MetricsServerParams struct {
	fx.In

	Srv *http.Server   `name:"fxmetrics"`
	Mux *http.ServeMux `name:"fxmetrics"`
	Lc  fx.Lifecycle
	Log *slog.Logger
}

func ServerInvoker(p MetricsServerParams) {
	p.Srv.Handler = p.Mux
	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := p.Srv.ListenAndServe()
				p.Log.Error("listening metrics", "addr", p.Srv.Addr)
				if errors.Is(err, http.ErrServerClosed) {
					return
				}
				if err != nil {
					p.Log.Error("failed listening metrics", "err", err)
				}
			}()
			p.Log.Info("debug api", "addr", p.Srv.Addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.Srv.Shutdown(ctx)
		},
	})
}
