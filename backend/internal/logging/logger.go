package logging

import (
	"io"
	"log/slog"
	"os"
)

type HandlerType string

const (
	TextHandler HandlerType = "text"
	JSONHandler HandlerType = "json"
)

// NewLogger creates a slog.Logger with the given level and handler type.
// Output is written to writer; if writer is nil, os.Stdout is used.
func NewLogger(level slog.Level, handlerType HandlerType, writer io.Writer) *slog.Logger {
	if writer == nil {
		writer = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	switch handlerType {
	case JSONHandler:
		handler = slog.NewJSONHandler(writer, opts)
	default:
		handler = slog.NewTextHandler(writer, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
