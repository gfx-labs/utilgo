package gotel

import (
	"context"
	"gfx.cafe/util/go/fxplus"
	"go.uber.org/fx"
	"log/slog"
)

type TraceProvider struct {
	Enabled     bool
	ServiceName string
}

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
	o.ServiceName = string(p.ServiceName)

	f, err := InitTracing(context.Background(), WithServiceName(o.ServiceName))
	if err == nil {
		o.Enabled = true

		p.Lc.Append(fx.Hook{
			OnStop: f,
		})
	} else {
		p.Log.Warn("error initializing tracing", "err", err)
	}

	r.Output = o
	return
}
