// Package webkit provides a lightweight HTTP framework built on top of Chi router.
//
// It offers a simple API for routing, middleware (called "plugs"), error handling,
// and context management. The framework supports graceful shutdown and provides
// utilities for common HTTP operations like JSON responses, file serving, and redirects.
//
// Basic usage:
//
//	app := webkit.New()
//	app.Plug(middleware.Logger)
//	app.Get("/", func(c *webkit.Context) error {
//	    return c.JSON(http.StatusOK, map[string]string{"message": "hello"})
//	})
//	app.Start(":8080")
package webkit
