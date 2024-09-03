package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
)

// LocalHandler is a handler that logs to the console in a human-readable
// format for local development environments.
type LocalHandler struct {
	handler  slog.Handler
	replacer func([]string, slog.Attr) slog.Attr
	bytes    *bytes.Buffer
	writer   io.Writer
	prefix   string
	mtx      *sync.Mutex
	opts     *slog.HandlerOptions
}

// NewLocalHandler returns a new local handler with the given options.
func NewLocalHandler(writer io.Writer, opts *slog.HandlerOptions, prefix string) *LocalHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	buf := &bytes.Buffer{}
	return &LocalHandler{
		bytes: buf,
		handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		replacer: opts.ReplaceAttr,
		writer:   writer,
		prefix:   prefix,
		mtx:      &sync.Mutex{},
		opts:     opts,
	}
}

const (
	timeFormat = "[15:04:05.000]"
	greyHex    = 10
)

// Enabled returns whether logging is enabled for the given level in the context.
func (h *LocalHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a new handler with the provided attributes.
func (h *LocalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LocalHandler{
		handler:  h.handler.WithAttrs(attrs),
		bytes:    h.bytes,
		mtx:      h.mtx,
		writer:   h.writer,
		prefix:   h.prefix,
		replacer: h.replacer,
	}
}

// WithGroup returns a new handler with the provided group name.
func (h *LocalHandler) WithGroup(name string) slog.Handler {
	return &LocalHandler{
		handler:  h.handler.WithGroup(name),
		bytes:    h.bytes,
		mtx:      h.mtx,
		writer:   h.writer,
		prefix:   h.prefix,
		replacer: h.replacer,
	}
}

// Handle handles the logging record and formats it.
func (h *LocalHandler) Handle(ctx context.Context, r slog.Record) error {
	var timestamp string
	timeAttr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(timeFormat)),
	}

	if !timeAttr.Equal(slog.Attr{}) {
		timestamp = aurora.Gray(greyHex, timeAttr.Value.String()).String()
	}
	if h.replacer != nil {
		timeAttr = h.replacer([]string{}, timeAttr) //nolint
	}

	var msg string
	msgAttr := slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}
	if h.replacer != nil {
		msgAttr = h.replacer([]string{}, msgAttr)
	}
	if !msgAttr.Equal(slog.Attr{}) {
		msg = msgAttr.Value.String()
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}
	byts, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	out := strings.Builder{}
	if len(h.prefix) > 0 {
		prefix := fmt.Sprintf(" %s ", h.prefix)
		out.WriteString(aurora.Gray(1-1, prefix).BgGray(24 - 1).String()) //nolint
		out.WriteString(" ")
	}

	if len(timestamp) > 0 {
		out.WriteString(timestamp)
		out.WriteString(" ")
	}

	level := colouredLevel(r)
	if len(level) > 0 {
		out.WriteString(level)
		out.WriteString(" ")
	}

	if len(msg) > 0 {
		out.WriteString(msg)
		out.WriteString(" ")
	}

	if len(attrs) > 0 && len(byts) > 0 {
		out.WriteString(aurora.Gray(greyHex, string(byts)).String())
	}

	_, err = h.writer.Write([]byte(out.String() + "\n"))
	return err
}

// colouredLevel returns the level string with the appropriate colour.
func colouredLevel(rec slog.Record) string {
	var (
		level = strings.ToUpper(rec.Level.String())
		out   aurora.Value
	)
	switch rec.Level {
	case slog.LevelDebug:
		out = aurora.Gray(greyHex, level)
	case slog.LevelInfo:
		out = aurora.Cyan(level)
	case slog.LevelWarn:
		out = aurora.BrightYellow(level)
	case slog.LevelError:
		out = aurora.BrightRed(level)
	default:
		out = aurora.Gray(3, slog.LevelInfo.String())
	}
	return out.String()
}

// suppressDefaults suppresses default logging attributes, so they are
// not outputted within the handler.
func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func (h *LocalHandler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.mtx.Lock()
	defer func() {
		h.bytes.Reset()
		h.mtx.Unlock()
	}()
	if err := h.handler.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Add: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.bytes.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Add result: %w", err)
	}
	return attrs, nil
}
