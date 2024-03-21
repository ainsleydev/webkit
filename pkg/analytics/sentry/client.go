package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/ainsleydev/webkit/pkg/env"
)

// Init initialises the Sentry client and returns a function to close it.
func Init(dsn string) (func(), error) {
	if env.IsDevelopment() {
		return func() {}, nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            false,
		AttachStacktrace: true,
		Environment:      env.AppEnvironment(),
	})
	if err != nil {
		return func() {}, err
	}

	// Example Message
	// sentry.CaptureMessage("It works!")

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	return func() {
		sentry.Flush(2 * time.Second)
	}, nil
}
