/*
Copyright (c) 2015-present Peter Kieltyka (https://github.com/pkieltyka), Google Inc.
Modified by: 2023 Eli Janssen

# MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

~~~~~~

Originally from:
  https://github.com/go-chi/chi/blob/58ca6d6119ed77f8c1d564bac109fc36db10e3d0/middleware/request_id.go
*/

package header

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/dropwhile/refid/v2"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/util"
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

var (
	prefix string
	reqid  uint64
)

func init() {
	refID := util.Must(refid.New())
	prefix = refID.String()
}

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			myid := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%s-%06d", prefix, myid)
		}
		ctx = context.WithValue(ctx, middleware.RequestIDKey, requestID)
		ctx = logger.PrependAttr(ctx, "request_id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

/*
// NextRequestID generates the next request ID in the sequence.
func NextRequestID() uint64 {
	return atomic.AddUint64(&reqid, 1)
}

// type AttrExtractor func(ctx context.Context, recordT time.Time, recordLvl slog.Level, recordMsg string) []slog.Attr
func RequestIDExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	reqID := GetReqID(ctx)
	if reqID != "" {
		return []slog.Attr{slog.String("request_id", reqID)}
	}
	return nil
}
*/
