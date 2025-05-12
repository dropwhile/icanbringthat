// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package debug

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/httplog/v2"
)

var skipPaths = []string{}

var logOptions = httplog.Options{
	Concise:        true,
	RequestHeaders: true,
	HideRequestHeaders: []string{
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
	},
	QuietDownRoutes: []string{
		"/",
		"/ping",
	},
	QuietDownPeriod: 10 * time.Second,
}

func RequestLogger() func(next http.Handler) http.Handler {
	httpLogger := &httplog.Logger{
		Logger:  slog.Default(),
		Options: logOptions,
	}
	return httplog.Handler(httpLogger, skipPaths)
}
