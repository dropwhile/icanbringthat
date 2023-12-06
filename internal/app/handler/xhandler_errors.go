package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
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
	log.Debug().CallerSkipFrame(1).Err(err).Msgf("error parsing form: %s", strings.Join(hints, ", "))
	errStr := "bad form data"
	if len(hints) > 0 {
		errStr = fmt.Sprintf("%s - %s", errStr, strings.Join(hints, ", "))
	}
	http.Error(w, errStr, http.StatusBadRequest)
}

/* not found */

func (x *Handler) BadRefIDError(w http.ResponseWriter, reftyp string, err error) {
	log.Debug().CallerSkipFrame(1).Err(err).Msgf("bad %s ref-id", reftyp)
	x.Error(w, fmt.Sprintf("Invalid %s-ref-id", reftyp), http.StatusNotFound)
}

func (x *Handler) NotFoundError(w http.ResponseWriter) {
	x.Error(w, "Not Found", http.StatusNotFound)
}

/* internal server error */

func (x *Handler) DBError(w http.ResponseWriter, err error) {
	log.Info().CallerSkipFrame(1).Stack().Err(err).Msg("db error")
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
		log.Debug().CallerSkipFrame(1).Err(err).Msg("custom error status page render issue")
		fmt.Fprintln(w, statusMsg)
		return
	}
}
