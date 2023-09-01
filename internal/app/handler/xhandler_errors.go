// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dropwhile/icanbringthat/internal/logger"
)

/*** specific errors ***/

/* forbidden */

func (x *Handler) AccessDeniedError(w http.ResponseWriter) {
	x.Error(w, "Access Denied", http.StatusForbidden)
}

/* bad request */

func (x *Handler) BadSessionDataError(w http.ResponseWriter) {
	x.Error(w, "Bad Session Data", http.StatusBadRequest)
}

func (x *Handler) BadFormDataError(w http.ResponseWriter, err error, hints ...string) {
	logger.LogSkip(slog.Default(), 1, slog.LevelDebug,
		context.Background(),
		"error parsing form",
		"hints", strings.Join(hints, ", "), "error", err)
	errStr := "bad form data"
	if len(hints) > 0 {
		errStr = fmt.Sprintf("%s - %s", errStr, strings.Join(hints, ", "))
	}
	http.Error(w, errStr, http.StatusBadRequest)
}

/* not found */

func (x *Handler) BadRefIDError(w http.ResponseWriter, reftyp string, err error) {
	logger.LogSkip(slog.Default(), 1, slog.LevelDebug,
		context.Background(),
		"bad ref-id", "type", reftyp, "error", err)
	x.Error(w, fmt.Sprintf("Invalid %s-ref-id", reftyp), http.StatusNotFound)
}

func (x *Handler) NotFoundError(w http.ResponseWriter) {
	x.Error(w, "Not Found", http.StatusNotFound)
}

/* internal server error */

func (x *Handler) DBError(w http.ResponseWriter, err error) {
	logger.LogSkip(slog.Default(), 1, slog.LevelDebug,
		context.Background(),
		"db error", "error", err)
	x.Error(w, "DB error", http.StatusInternalServerError)
}

func (x *Handler) TemplateError(w http.ResponseWriter) {
	x.Error(w, "Template Render Error", http.StatusInternalServerError)
}

/*** general errors ***/

func (x *Handler) BadRequestError(w http.ResponseWriter, msg string) {
	x.Error(w, msg, http.StatusBadRequest)
}

func (x *Handler) InternalServerError(w http.ResponseWriter, msg string) {
	x.Error(w, msg, http.StatusInternalServerError)
}

func (x *Handler) ForbiddenError(w http.ResponseWriter, msg string) {
	x.Error(w, msg, http.StatusForbidden)
}

func (x *Handler) Error(w http.ResponseWriter, statusMsg string, code int) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	err := x.TemplateExecute(w, "error-page.gohtml", MapSA{
		"ErrorCode":   code,
		"ErrorStatus": statusMsg,
		"title":       fmt.Sprintf("%d - %s", code, statusMsg),
	})
	if err != nil {
		// error rendering template, so just return a very basic status page
		logger.LogSkip(slog.Default(), 1, slog.LevelDebug,
			context.Background(),
			"custom error status page render issue", "error", err)
		fmt.Fprintln(w, statusMsg)
		return
	}
}
