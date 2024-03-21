package app

import (
	"github.com/ainsleydev/webkit/pkg/analytics/sentry"
	log "github.com/ainsleydev/webkit/pkg/logger"
)

func Bootstrap() {

	log.Bootstrap("Prefix")

	sentry.Init("DSN")
}
