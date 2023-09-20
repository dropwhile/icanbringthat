package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ShowCreateEventItemForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event")
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	tplVars := map[string]any{
		"user":           user,
		"event":          event,
		"title":          "Create Event Item",
		"nav":            "create-event-item",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if Hx(r).Target() == "modalbody" {
		err = h.TemplateExecuteSub(w, "create-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = h.TemplateExecute(w, "create-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowEventItemEditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		http.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event")
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	eventItem, err := model.GetEventItemByRefId(ctx, h.Db, eventItemRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event item")
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	tplVars := map[string]any{
		"user":           user,
		"event":          event,
		"eventItem":      eventItem,
		"title":          "Edit Event Item",
		"nav":            "edit-event-item",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if Hx(r).Target() == "modalbody" {
		err = h.TemplateExecuteSub(w, "edit-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = h.TemplateExecute(w, "edit-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		log.Debug().Err(err).Msg("bad event ref-id")
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for event")
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		log.Debug().Err(err).Msg("description form element empty")
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	_, err = model.NewEventItem(ctx, h.Db, event.Id, description)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefId), http.StatusSeeOther)
}

func (h *Handler) UpdateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		http.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	eventItem, err := model.GetEventItemByRefId(ctx, h.Db, eventItemRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if eventItem.EventId != event.Id {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.Id", event.Id).
			Int("eventItem.EventId", eventItem.EventId).
			Msg("eventItem.EventId and event.Id mismatch")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// check if earmark exists, and is marked by someone else
	// if so, disallow editing in that case.
	earmark, err := model.GetEarmarkByEventItem(ctx, h.Db, eventItem)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	case err == nil:
		if earmark.UserId != user.Id {
			log.Info().
				Int("user.Id", user.Id).
				Int("earmark.UserId", earmark.UserId).
				Msg("user id mismatch")
			http.Error(w, "earmarked by other user - access denied", http.StatusForbidden)
			return
		}
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		log.Debug().Err(err).Msg("missing form data")
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	eventItem.Description = description
	err = eventItem.Save(ctx, h.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", eventRefId)) {
		w.Header().Add("HX-Refresh", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		http.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	eventItem, err := model.GetEventItemByRefId(ctx, h.Db, eventItemRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if eventItem.EventId != event.Id {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.Id", event.Id).
			Int("eventItem.EventId", eventItem.EventId).
			Msg("eventItem.EventId and event.Id mismatch")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err = eventItem.Delete(ctx, h.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
