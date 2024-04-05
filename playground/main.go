package main

import (
	"log/slog"
	"os"

	"github.com/ainsleydev/webkit/pkg/app"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/log"
	"github.com/ainsleydev/webkit/pkg/middleware"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

type Env struct {
}

func main() {
	kit := webkit.New()

	log.Bootstrap("Playground")

	err := env.ParseConfig(&Env{})
	if err != nil {
		slog.Error("Failed to parse config: %v", err)
		os.Exit(1)
	}

	app.StartInternalHTTP()

	kit.Plug(middleware.Logger)
	kit.Plug(middleware.Recover)
	kit.Plug(middleware.TrailingSlashRedirect)
	kit.Plug(middleware.NonWWWRedirect)
	kit.Plug(middleware.RequestID)
	kit.Plug(middleware.Gzip)
	kit.Plug(middleware.CORS)

	kit.Get("/ping", webkit.PingHandler)
	kit.Get("/", func(ctx *webkit.Context) error {
		return ctx.String(500, "Hello, Crab Poo!")
	})

	if err := kit.Start(":8080"); err != nil {
		slog.Error("Failed to start server: %v", err)
	}
}
