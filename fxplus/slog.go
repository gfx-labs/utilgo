package fxplus

import (
	"log/slog"
	"os"

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
	return &fxevent.SlogLogger{Logger: logger}
})
