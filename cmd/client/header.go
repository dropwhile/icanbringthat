// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
)

func setHeaders(into, from http.Header, merge bool) {
	for key, vals := range from {
		if len(vals) == 0 {
			// For response trailers, net/http will pre-populate entries
			// with nil values based on the "Trailer" header. But if there
			// are no actual values for those keys, we skip them.
			continue
		}
		if merge {
			into[key] = append(into[key], vals...)
		} else {
			into[key] = vals
		}
	}
}

func NewAddHeadersInterceptor(headers http.Header) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			setHeaders(req.Header(), headers, false)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
