package gotel

import (
	"context"
	"gfx.cafe/util/go/fxplus"
	"go.uber.org/fx"
	"log/slog"
)

type TraceProvider struct {
	enabled     bool
	serviceName string
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
	o.serviceName = string(p.ServiceName)

	f, err := InitTracing(context.Background(), WithServiceName(o.serviceName))
	if err == nil {
		o.enabled = true

		p.Lc.Append(fx.Hook{
			OnStop: f,
		})
	} else {
		p.Log.Warn("error initializing tracing", "err", err)
	}

	r.Output = o
	return
}

func (p *TraceProvider) Enabled() bool {
	if p == nil {
		return false
	}

	return p.enabled
}

func (p *TraceProvider) ServiceName() string {
	if p == nil {
		return ""
	}

	return p.serviceName
}
