package logger

import (
	"log/slog"
	"os"
)

func New(environment string) *slog.Logger {
	if environment == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
