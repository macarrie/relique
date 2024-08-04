package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func Init(debug bool) {
	var logLevel slog.Level
	if debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	opts := &tint.Options{
		Level: logLevel,
	}

	logger := slog.New(tint.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}
