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

// SentryHook implements logrus.Hook to send errors to sentry
type SentryHook struct {
	client *sentry.Client
}

// NewSentryHook creates a sentry hook for logrus given a sentry client
func NewSentryHook(client *sentry.Client) SentryHook {
	return SentryHook{
		client: client,
	}
}

// Levels returns the levels this hook is enabled for. This is a part
// of logrus.Hook.
func (h SentryHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire is an event handler for logrus. This is a part of logrus.Hook.
// Taken from: https://github.com/getsentry/sentry-go/issues/43
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

// SentryEventIdentityModifier is a sentry event modifier that simply passes
// through the event.
type SentryEventIdentityModifier struct{}

// ApplyToEvent simply returns the event (ignoring the hint).
func (m *SentryEventIdentityModifier) ApplyToEvent(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	return event
}

var sentryModifier = &SentryEventIdentityModifier{}
