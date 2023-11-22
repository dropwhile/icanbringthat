package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
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
		x.Error(w, "bad webhook data error", http.StatusBadRequest)
		return
	}

	if pm.RecordType != "SubscriptionChange" {
		log.Info().Any("postmark", pm).Msg("unexpecte RecordType")
		w.WriteHeader(http.StatusOK)
		return
	}

	user, err := model.GetUserByEmail(ctx, x.Db, pm.Recipient)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no user found from callback")
		w.WriteHeader(http.StatusOK)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// we only disable reminders, not re-enable them
	if !pm.SuppressSending {
		w.WriteHeader(http.StatusOK)
		return
	}

	// if already disabled, no need to disable again
	if !user.Settings.EnableReminders {
		log.Info().
			Any("postmark", pm).
			Msg("reminders already disabled")
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Info().
		Any("postmark", pm).
		Msg("disabling reminders due to postmark callback")

	// bounced email, marked spam, unsubscribed...etc
	// so... disable reminders
	user.Settings.EnableReminders = false
	err = pgx.BeginFunc(ctx, x.Db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUserSettings(ctx, tx, &user.Settings, user.ID)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = model.NewNotification(ctx, tx, user.ID,
			fmt.Sprintf("email notifications disabled due to '%s'", pm.SuppressionReason),
		)
		if innerErr != nil {
			return innerErr
		}
		return nil
	})
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}
}
