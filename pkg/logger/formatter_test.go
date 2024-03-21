package log

import (
	"bytes"
	"log/slog"
	"sync"
	"testing"

	"github.com/logrusorgru/aurora"
	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input func()
		want  string
	}{
		"Info level with message": {
			input: func() {
				Debug(" Test")
			},
			want: "DEBUG test message",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			slog.SetDefault(slog.New(&LocalHandler{
				handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
				bytes:   buf,
				prefix:  "test_prefix",
				writer:  buf,
				mtx:     &sync.Mutex{},
			}))
			test.input()
			assert.Contains(t, buf.String(), test.want, "Output should contain the expected message")
		})
	}
}

func TestColouredLevel(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		level slog.Level
		want  string
	}{
		"Debug": {
			level: slog.LevelDebug,
			want:  aurora.Gray(10, "DEBUG").String(),
		},
		"Info": {
			level: slog.LevelInfo,
			want:  aurora.Cyan("INFO").String(),
		},
		"Warning": {
			level: slog.LevelWarn,
			want:  aurora.BrightYellow("WARN").String(),
		},
		"Error": {
			level: slog.LevelError,
			want:  aurora.BrightRed("ERROR").String(),
		},
		"Default": {
			level: slog.Level(100),
			want:  aurora.Gray(3, "INFO").String(),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := colouredLevel(slog.Record{Level: test.level})
			assert.Equal(t, test.want, got, "Coloured level should match")
		})
	}
}
