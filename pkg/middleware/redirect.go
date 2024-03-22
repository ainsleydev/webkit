package middleware

// See: https://github.com/labstack/echo/blob/master/middleware/redirect.go

// HTTPSNonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSNonWWWRedirect()`.
//func HTTPSNonWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFunc {
//	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
//		if scheme != "https" {
//			host = strings.TrimPrefix(host, www)
//			return true, "https://" + host + uri
//		}
//		return false, ""
//	})
//}
