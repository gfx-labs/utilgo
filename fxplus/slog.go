package fxplus

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"github.com/gfx-labs/utilgo/whereami"
)

func NewLogger() *slog.Logger {
	slog_level := slog.LevelInfo.Level()
	if GO_LOG := os.Getenv("GO_LOG"); GO_LOG != "" {
		slog_level = ParseSlogFilter(slog.LevelInfo, GO_LOG)
	}

	// Use JSON handler in Docker or Kubernetes environments
	if whereami.IsDocker() || whereami.IsKubernetes() {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog_level,
		}))
	}

	// Use tint handler for local development
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource: true,
		Level:     slog_level,
	}))
}

var WithLogger = fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
	l := &fxevent.SlogLogger{Logger: logger}
	l.UseLogLevel(slog.LevelDebug)
	return l
})

func ParseSlogFilter(defaultLevel slog.Level, filter string) slog.Level {
	filters := strings.Split(filter, ",")
	for _, filter := range filters {
		// TODO: think about what to do with this :)
		first, _, ok := strings.Cut(filter, "=")
		if !ok {
			defaultLevel.UnmarshalText([]byte(first))
			continue
		}
	}

	return defaultLevel
}
