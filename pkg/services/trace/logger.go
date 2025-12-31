package trace

import (
	"log/slog"
	"os"
)

type Kind string

const (
	JSON Kind = "json"
	TEXT Kind = "text"
)

func Logger(kind Kind, debug bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	writer := os.Stdout
	options := &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	}
	handler := slog.Handler(nil)

	switch kind {
	case JSON:
		handler = slog.NewJSONHandler(writer, options)
	case TEXT:
		handler = slog.NewTextHandler(writer, options)
	}

	return slog.New(handler)
}
