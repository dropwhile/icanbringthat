package debug

import (
	"net/http"
	"time"

	"github.com/go-chi/httplog"
)

var defaultIgnoreHeaders = []string{
	"accept",
	"accept-encoding",
	"accept-language",
	"accept-ranges",
	"connection",
	"cookie",
	"sec-ch-ua",
	"sec-ch-ua-mobile",
	"sec-ch-ua-platform",
	"sec-fetch-dest",
	"sec-fetch-mode",
	"sec-fetch-site",
	"sec-fetch-user",
	"sec-gpc",
	"upgrade-insecure-requests",
	"user-agent",
	"scheme",
	"x-csrf-token",
}

func RequestLogger() func(next http.Handler) http.Handler {
	return httplog.RequestLogger(
		httplog.NewLogger(
			"httplog-example",
			httplog.Options{
				SkipHeaders:     defaultIgnoreHeaders,
				TimeFieldFormat: time.RFC3339,
			},
		),
	)
}
