package gotel

import (
	"context"
	"github.com/gfx-labs/utilgo/fxplus"
	"go.uber.org/fx"
	"log/slog"
)

type TraceProvider struct{}

type Params struct {
	fx.In

	Ctx         context.Context
	Lc          fx.Lifecycle
	Log         *slog.Logger
	ServiceName fxplus.ComponentName
}

type Result struct {
	fx.Out

	Output *TraceProvider
}

func NewTraceProvider(p Params) (r Result, err error) {
	o := &TraceProvider{}

	f, err := InitTracing(context.Background(), WithServiceName(string(p.ServiceName)))
	if err == nil {
		p.Lc.Append(fx.Hook{
			OnStop: f,
		})
	} else {
		p.Log.Warn("error initializing tracing", "err", err)
	}

	r.Output = o
	return
}
