package app

import (
	"github.com/ainsleydev/webkit/pkg/analytics/sentry"
	log "github.com/ainsleydev/webkit/pkg/log"
)

func Bootstrap() {

	// Unmarshal?

	log.Bootstrap("Prefix")

	sentry.Init("DSN")
}
