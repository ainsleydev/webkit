package main

import (
	"log/slog"
	"os"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/log"
	"github.com/ainsleydev/webkit/pkg/middleware"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

type Env struct {
}

func main() {
	app := webkit.New()

	log.Bootstrap("Playground")

	err := env.ParseConfig(&Env{})
	if err != nil {
		slog.Error("Failed to parse config: %v", err)
		os.Exit(1)
	}

	app.Plug(middleware.Logger)
	app.Plug(middleware.Recover)
	app.Plug(middleware.RedirectSlashes)
	app.Plug(middleware.RequestID)
	app.Plug(middleware.Gzip)

	app.Get("/", func(ctx *webkit.Context) error {
		return ctx.String(500, "Hello, World!")
	})

	if err := app.Start(":8080"); err != nil {
		slog.Error("Failed to start server: %v", err)
	}
}
