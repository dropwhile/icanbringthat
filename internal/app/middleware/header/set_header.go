package header

import (
	"net/http"
)

// SetHeader is a convenience handler to set a response header key/value
func SetHeader(key, value string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func SetHeaders(headers map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
