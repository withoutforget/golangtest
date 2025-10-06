package logging

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
)

func InitLogging() (*slog.Logger, error) {
	logger := slog.New(console.NewHandler(
		os.Stdout,
		&console.HandlerOptions{Level: slog.LevelDebug, AddSource: true},
	))

	slog.SetDefault(logger)

	return slog.Default(), nil
}
