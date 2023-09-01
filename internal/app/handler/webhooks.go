// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"encoding/json"
	"log/slog"
	"mime"
	"net/http"
	"strings"

	"github.com/dropwhile/icanbringthat/internal/errs"
)

type PostMarkRecord struct {
	RecordType        string
	MessageStream     string
	SuppressionReason string
	Recipient         string
	SuppressSending   bool
}

func (x *Handler) PostmarkCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isJson := false
	for _, v := range strings.Split(r.Header.Get("Content-type"), ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == "application/json" {
			isJson = true
			break
		}
	}
	if !isJson {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	var pm PostMarkRecord
	err := dec.Decode(&pm)
	if err != nil {
		slog.InfoContext(ctx, "webhook error", "error", err)
		x.BadRequestError(w, "bad webhook data")
		return
	}

	if pm.RecordType != "SubscriptionChange" {
		slog.InfoContext(ctx, "unexpecte RecordType", "postmark", pm)
		w.WriteHeader(http.StatusOK)
		return
	}

	// we only disable reminders, not re-enable them
	if !pm.SuppressSending {
		w.WriteHeader(http.StatusOK)
		return
	}

	slog.DebugContext(ctx,
		"disabling reminders due to postmark callback",
		"postmark", pm)
	errx := x.svc.DisableRemindersWithNotification(
		ctx, pm.Recipient, pm.SuppressionReason,
	)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			slog.InfoContext(ctx, "no user found from callback",
				"error", err)
			w.WriteHeader(http.StatusOK)
		case errs.FailedPrecondition:
			slog.InfoContext(ctx, "reminders already disabled",
				"postmark", pm)
			w.WriteHeader(http.StatusOK)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}
}
