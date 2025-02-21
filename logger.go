package servicekit

import (
	"context"
	"log/slog"
	"os"
)

type key string

const (
	fields key = "slog_fields"
)

type Handler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h Handler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(fields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(fields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, fields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, fields, v)
}

type LoggerConfig interface {
	LoggerLevel() string
	LoggerSource() bool
	LoggerFormat() string
}

func SetLogger(config LoggerConfig) {
	var handler slog.Handler

	var level slog.Level
	switch config.LoggerLevel() {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		slog.Error("Invalid log level", "level", config.LoggerLevel)
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		AddSource: config.LoggerSource(),
		Level:     level,
	}

	switch config.LoggerFormat() {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		slog.Error("Invalid log format", "format", config.LoggerFormat())
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(&Handler{Handler: handler})
	slog.SetDefault(logger)
}
