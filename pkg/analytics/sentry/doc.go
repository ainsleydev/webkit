// Package sentry provides Sentry error tracking integration for WebKit applications.
//
// It initializes the Sentry client with environment-aware configuration and
// automatically disables tracking in development. Returns a cleanup function
// to flush events before shutdown.
package sentry
