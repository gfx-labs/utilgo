package fxplus

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func NewLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource: true,
		Level:     slog.LevelInfo.Level(),
	}))
}

var WithLogger = fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
	return &eventLogger{L: logger}
})

type eventLogger struct {
	L *slog.Logger
}

func (l *eventLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.L.Debug("[Fx] OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.L.Error("[Fx] OnStart hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"err", e.Err,
			)
		} else {
			l.L.Debug("[Fx] OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		l.L.Debug("[Fx] OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.L.Error("[Fx] OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"err", e.Err,
			)
		} else {
			l.L.Debug("[Fx] OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.L.Error("[Fx] error encountered while applying options",
				"type", e.TypeName,
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
				"err", e.Err,
			)
		} else {
			l.L.Debug("[Fx] supplied",
				"type", e.TypeName,
				//	"stacktrace", e.StackTrace,
				//			"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.L.Debug("[Fx] provided",
				"constructor", e.ConstructorName,
				//		"stacktrace", e.StackTrace,
				//				"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.L.Error("[Fx] error encountered while applying options",
				"module", e.ModuleName,
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"err", e.Err,
			)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.L.Debug("[Fx] replaced",
				//			"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.L.Error("[Fx] error encountered while replacing",
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
				"err", e.Err,
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.L.Debug("[Fx] decorated",
				"decorator", e.DecoratorName,
				//			"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.L.Error("[Fx] error encountered while applying options",
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"module", e.ModuleName,
				"err", e.Err,
			)
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.L.Error("[Fx] error returned",
				"name", e.Name,
				"kind", e.Kind,
				"module", e.ModuleName,
				"err", e.Err,
			)
		} else {
			l.L.Debug("[Fx] run",
				"name", e.Name,
				"kind", e.Kind,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.L.Debug("[Fx] invoking",
			"function", e.FunctionName,
			"module", e.ModuleName,
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.L.Error("[Fx] invoke failed",
				"err", e.Err,
				"stack", e.Trace,
				"function", e.FunctionName,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Stopping:
		l.L.Debug("[Fx] received signal",
			"signal", strings.ToUpper(e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			l.L.Error("[Fx] stop failed", "err", e.Err)
		}
	case *fxevent.RollingBack:
		l.L.Error("[Fx] start failed, rolling back", "err", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.L.Error("[Fx] rollback failed", "err", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.L.Error("[Fx] start failed", "err", e.Err)
		} else {
			l.L.Info("[Fx] app started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.L.Error("[Fx] custom logger initialization failed", "err", e.Err)
		} else {
			l.L.Debug("[Fx] initialized custom fxevent.Logger", "function", e.ConstructorName)
		}
	}
}
