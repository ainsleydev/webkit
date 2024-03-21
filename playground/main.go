package main

import (
	"os"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/logger"
	"github.com/ainsleydev/webkit/pkg/middleware"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

type Env struct {
}

func main() {
	app := webkit.New()

	logger.Bootstrap("Playground")

	err := env.ParseConfig(&Env{})
	if err != nil {
		logger.Error("Failed to parse config: %v", err)
		os.Exit(1)
	}

	app.Plug(middleware.Logger)
	app.Plug(middleware.RequestID)
	app.Plug(middleware.RedirectSlashes)

	app.Get("/", func(ctx *webkit.Context) error {
		return ctx.String(200, "Hello, World!")
	})

	if err := app.Start(":8080"); err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}
