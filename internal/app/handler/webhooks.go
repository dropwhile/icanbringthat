package handler

import (
	"encoding/json"
	"mime"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/somerr"
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
		log.Info().Err(err).Msg("webhook error")
		x.BadRequestError(w, "bad webhook data")
		return
	}

	if pm.RecordType != "SubscriptionChange" {
		log.Info().Any("postmark", pm).Msg("unexpecte RecordType")
		w.WriteHeader(http.StatusOK)
		return
	}

	// we only disable reminders, not re-enable them
	if !pm.SuppressSending {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Info().
		Any("postmark", pm).
		Msg("disabling reminders due to postmark callback")
	errx := service.DisableRemindersWithNotification(
		ctx, x.Db, pm.Recipient, pm.SuppressionReason,
	)
	if errx != nil {
		switch errx.Code() {
		case somerr.NotFound:
			log.Info().Err(err).Msg("no user found from callback")
			w.WriteHeader(http.StatusOK)
		case somerr.FailedPrecondition:
			log.Info().
				Any("postmark", pm).
				Msg("reminders already disabled")
			w.WriteHeader(http.StatusOK)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}
}
