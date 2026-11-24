package log

import (
	"errors"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSentryHook(t *testing.T) {
	t.Parallel()

	client, tr := newTestSentryClient(t)
	hook := NewSentryHook(client)

	assert.Equal(t, []logrus.Level{logrus.FatalLevel, logrus.PanicLevel}, hook.Levels())

	err := hook.Fire(&logrus.Entry{
		Message: "crashed",
		Level:   logrus.FatalLevel,
		Time:    time.Now(),
		Data: logrus.Fields{
			logrus.ErrorKey: errors.New("boom"),
			"service":       "api",
		},
	})
	require.NoError(t, err)

	events := tr.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "crashed", events[0].Message)
	assert.Equal(t, sentry.LevelFatal, events[0].Level)
	assert.Equal(t, "api", events[0].Extra["service"])
	require.Len(t, events[0].Exception, 1)
	assert.Equal(t, "boom", events[0].Exception[0].Value)
	assert.Equal(t, 1, tr.flushes)
}
