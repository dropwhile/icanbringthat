package strip

import "net/http"

// StripPrefix is a middleware that will strip the provided prefix from the
// request path before handing the request over to the next handler.
func StripPrefix(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}
