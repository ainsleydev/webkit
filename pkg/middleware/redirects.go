package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// TrailingSlashRedirect is a middleware that will match request paths with a trailing
// slash and redirect to the same path, less the trailing slash.
//
// NOTE: RedirectSlashes middleware is *incompatible* with http.FileServer.
func TrailingSlashRedirect(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		r := ctx.Request
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] != '/' {
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s/?%s", path, r.URL.RawQuery)
			} else {
				path = path + "/"
			}
			redirectURL := fmt.Sprintf("//%s%s", r.Host, path)
			return ctx.Redirect(301, redirectURL)
		}
		return next(ctx)
	}
}

// redirectLogic represents a function that given a scheme, host and uri
// can both: 1) determine if redirect is needed (will set ok accordingly) and
// 2) return the appropriate redirect url.
type redirectLogic func(scheme, host, uri string) (ok bool, url string)

const www = "www."

// WWWRedirect redirects all requests to the www subdomain of the current host.
// Redirects are performed using a 301 status code (Permanent Redirect)
// using HTTPS.
//
// For example, a request to "https://example.com" will be redirected to "https://www.example.com".
func WWWRedirect(next webkit.Handler) webkit.Handler {
	return redirect(next, func(scheme, host, uri string) (bool, string) {
		if scheme != "https" && !strings.HasPrefix(host, www) {
			return true, "https://www." + host + uri
		}
		return false, ""
	})
}

// NonWWWRedirect redirects all requests to the non-www subdomain of the current host.
// Redirects are performed using a 301 status code (Permanent Redirect)
// using HTTPS.
//
// For example, a request to "https://www.example.com" will be redirected to "https://example.com".
func NonWWWRedirect(next webkit.Handler) webkit.Handler {
	return redirect(next, func(scheme, host, uri string) (bool, string) {
		if strings.HasPrefix(host, www) {
			return true, "https://" + host[4:] + uri
		}
		return false, ""
	})
}

// redirect returns a middleware that performs a redirect to the given URL.
func redirect(next webkit.Handler, cb redirectLogic) webkit.Handler {
	return func(ctx *webkit.Context) error {
		req := ctx.Request
		if ok, url := cb(ctx.Scheme(), req.Host, req.RequestURI); ok {
			return ctx.Redirect(http.StatusMovedPermanently, url)
		}
		return next(ctx)
	}
}
