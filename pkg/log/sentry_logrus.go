package log

import (
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

// severityMap is a mapping of logrus log level to sentry log level.
var severityMap = map[logrus.Level]sentry.Level{
	logrus.DebugLevel: sentry.LevelDebug,
	logrus.InfoLevel:  sentry.LevelInfo,
	logrus.WarnLevel:  sentry.LevelWarning,
	logrus.ErrorLevel: sentry.LevelError,
	logrus.FatalLevel: sentry.LevelFatal,
	logrus.PanicLevel: sentry.LevelFatal,
}

// SentryHook captures fatal and panic logs and sends them to Sentry for error tracking.
//
// Deprecated: SentryHook is a logrus hook and keeps logrus as a dependency
// of this package. Use NewSentryHandler, which works with slog. SentryHook
// will be removed in a future release.
type SentryHook struct {
	client *sentry.Client
}

// NewSentryHook creates a logrus hook that reports fatal and panic events to Sentry.
//
// Deprecated: Use NewSentryHandler with LevelFatal instead:
//
//	log.NewSentryHandler(next, client, &log.SentryHandlerOptions{Level: log.LevelFatal})
func NewSentryHook(client *sentry.Client) SentryHook {
	return SentryHook{
		client: client,
	}
}

// Levels returns the log levels this hook captures (fatal and panic only).
func (h SentryHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire captures the log entry and sends it to Sentry as an error event.
func (h SentryHook) Fire(e *logrus.Entry) error {
	event := sentry.NewEvent()
	event.Message = e.Message
	event.Level = severityMap[e.Level]
	event.Timestamp = e.Time

	var err error
	for k, v := range e.Data {
		if k == logrus.ErrorKey {
			err = v.(error)
		} else {
			event.Extra[k] = v
		}
	}

	if err != nil {
		stacktrace := sentry.ExtractStacktrace(err)
		event.Exception = []sentry.Exception{{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: stacktrace,
		}}
	}

	h.client.CaptureEvent(event, nil, sentryModifier)

	// Wait until the client has flushed all events to Sentry.
	// It is safe to wait a few seconds since this should only
	// event be called on a fatal or panic which will perform
	// and os.Exit.
	h.client.Flush(time.Second * 5)

	return nil
}
