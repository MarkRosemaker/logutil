package logutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	dir   = "logs"
	level = slog.LevelDebug
)

// Start starts logging, both in a file under logs/ named after the current time, and to stdout.
func Start(ctx context.Context) (func() error, error) {
	slog.SetLogLoggerLevel(level)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating log directory: %w", err)
	}

	path := filepath.Join(dir, time.Now().Format("20060102_150405")+".log")

	// Create the log file
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("creating log file: %w", err)
	}

	// Set the default logger to use a combined handler that writes to both the file and stdout
	slog.SetDefault(slog.New(slog.NewMultiHandler(
		slog.NewJSONHandler(file, &slog.HandlerOptions{Level: level}),
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	)))

	slog.InfoContext(ctx, "logging initialized", "path", path)

	return func() error {
		slog.InfoContext(ctx, "closing log file", "path", path)
		return file.Close()
	}, nil
}
