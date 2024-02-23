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
		fx.Provide(RouterHandler),
		fx.Invoke(ServerInvoker),
		fx.Options(opts...),
	)
}

type RouterParams struct {
	fx.In

	Routes []func(chi.Router) `group:"fxchi"`
}

func RouterHandler(p RouterParams) http.Handler {
	r := chi.NewRouter()
	for _, fn := range p.Routes {
		fn(r)
	}
	return r
}

func ServerInvoker(srv *http.Server, log *slog.Logger, lc fx.Lifecycle) {
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
