package middleware

import (
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See also [MDN: Cross-Origin Resource Sharing (CORS)].
//
// Security: Poorly configured CORS can compromise security because it allows
// relaxation of the browser's Same-Origin policy.  See [Exploiting CORS
// misconfigurations for Bitcoins and bounties] and [Portswigger: Cross-origin
// resource sharing (CORS)] for more details.
//
// [MDN: Cross-Origin Resource Sharing (CORS)]: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
// [Exploiting CORS misconfigurations for Bitcoins and bounties]: https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
// [Portswigger: Cross-origin resource sharing (CORS)]: https://portswigger.net/web-security/cors
func CORS(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {

		// See: https://github.com/labstack/echo/blob/master/middleware/cors.go

		//g.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//g.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//g.Writer.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, Content-Mime, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Origin, Cache-Control, X-Requested-With, token")
		//g.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		//
		//if g.Request.Method == "OPTIONS" {
		//	g.AbortWithStatus(http.StatusOK)
		//	return
		//}
		return next(ctx)
	}
}
