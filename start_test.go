package logutil_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/MarkRosemaker/logutil"
)

func TestStart(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Cleanup(func() { slog.SetDefault(slog.New(slog.DiscardHandler)) })

	closeLog, err := logutil.Start(t.Context(), slog.LevelInfo)
	if err != nil {
		t.Fatal(err)
	}

	slog.Debug("below the level")
	slog.Info("at the level")

	if err := closeLog(); err != nil {
		t.Fatal(err)
	}

	paths, err := filepath.Glob(filepath.Join("logs", "*.log"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("log files = %v, %v; want one", paths, err)
	}

	data, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}

	var msgs []string
	for line := range bytes.Lines(data) {
		var rec struct{ Msg string }
		if err := json.Unmarshal(line, &rec); err != nil {
			t.Fatalf("line %q: %v", line, err)
		}

		msgs = append(msgs, rec.Msg)
	}

	want := []string{"logging initialized", "at the level", "closing log file"}
	if len(msgs) != len(want) {
		t.Fatalf("logged %q, want %q", msgs, want)
	}

	for i := range want {
		if msgs[i] != want[i] {
			t.Fatalf("logged %q, want %q", msgs, want)
		}
	}
}
