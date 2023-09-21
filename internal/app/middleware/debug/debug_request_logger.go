package debug

import (
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
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
	httplog.DefaultOptions.SkipHeaders = defaultIgnoreHeaders
	return httplog.RequestLogger(log.Logger)
}
