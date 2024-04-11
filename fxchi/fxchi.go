package fxchi

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func HttpApp(opts ...fx.Option) fx.Option {
	return fx.Module("fxchi",
		fx.Provide(chi.NewRouter),
		fx.Provide(fx.Annotate(RouterHandler, fx.ResultTags(`name:"fxchi"`))),
		fx.Invoke(fx.Annotate(ServerInvoker, fx.ParamTags(`name:"fxchi"`, `name:"fxchi"`))),
		fx.Options(opts...),
	)
}

type RouterParams struct {
	fx.In

	Routes []func(chi.Router) `group:"fxchi"`
}

type RouteResults struct {
	fx.Out

	Route func(chi.Router) `group:"fxchi"`
}

func NewRoute(fn func(r chi.Router)) func() RouteResults {
	return func() RouteResults {
		return MakeRoute(fn)
	}
}

func MakeRoute(fn func(r chi.Router)) RouteResults {
	return RouteResults{
		Route: fn,
	}
}

func RouterHandler(r *chi.Mux, p RouterParams) http.Handler {
	for _, fn := range p.Routes {
		r.Group(fn)
	}
	return r
}

func ServerInvoker(srv *http.Server, handler http.Handler, log *slog.Logger, lc fx.Lifecycle) {
	srv.Handler = handler
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := srv.Addr
			if addr == "" {
				addr = ":http"
			}
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			go func() {
				log.Info("http server started", "addr", addr)
				srv.Serve(ln)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
