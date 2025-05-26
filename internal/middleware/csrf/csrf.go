// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// Inspired by
// * https://github.com/golang/go/issues/73626
// * https://github.com/tailscale/tailscale/pull/15735
package csrf

import (
	"net/http"
	"net/url"
)

type options struct {
	// AllowSecFetchSiteSameSite specifies whether to allow requests with the
	// Sec-Fetch-Site header set to "same-site" indicating that they are
	// cross-origin but that their origin shares the same site (gTLD+1) with
	// that of the request.
	AllowSecFetchSiteSameSite bool
}

type optionFunc func(*options)

func AllowSecFetchSiteSameSite() optionFunc {
	return func(o *options) {
		o.AllowSecFetchSiteSameSite = true
	}
}

// Protect routes against CSRF attacks by requiring non-(GET|HEAD|OPTIONS)
// requests to specify the Sec-Fetch-Site header with the value "same-origin",
// or if Sec-Fetch-Site is missing, with an Origin header matching the hostname
// in the Host header.
func Protect(optfuncs ...optionFunc) func(next http.Handler) http.Handler {
	opts := &options{}
	for _, optfunc := range optfuncs {
		optfunc(opts)
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Allow GET, HEAD, and OPTIONS requests without Sec-Fetch-Site
			// header checks.
			switch r.Method {
			case "GET", "HEAD", "OPTIONS":
				next.ServeHTTP(w, r)
				return
			}

			switch r.Header.Get("Sec-Fetch-Site") {
			case "same-origin":
				// allow same-origin requests
			case "same-site":
				// allow cross-origin, but same-site requests if configured
				if !opts.AllowSecFetchSiteSameSite {
					http.Error(w, "forbidden cross-origin request", http.StatusForbidden)
					return
				}
			case "cross-site", "none":
				// deny cross-site requests or direct navigation non-GET requests.
				http.Error(w, "forbidden cross-origin request", http.StatusForbidden)
				return
			default:
				// if Origin is present and is the same as Host, allow
				if origin := r.Header.Get("origin"); origin != "" {
					if u, err := url.Parse(origin); err == nil {
						if u.Host == r.Host {
							next.ServeHTTP(w, r)
							return
						}
					}
					http.Error(w, "forbidden cross-origin request", http.StatusForbidden)
					return
				}

				// neither sec-fetch-site, nor origin headers present
				// this is probably not a browser -- deny request
				http.Error(w,
					"missing required Sec-Fetch-Site header. You might need to update your browser.",
					http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
