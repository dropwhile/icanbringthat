package zhandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util/htmx"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (z *ZHandler) ShowCreateEventItemForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		z.Error(w, "access denied", http.StatusForbidden)
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
	if htmx.Hx(r).Target() == "modalbody" {
		err = z.TemplateExecuteSub(w, "create-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = z.TemplateExecute(w, "create-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) ShowEventItemEditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		z.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		z.Error(w, "access denied", http.StatusForbidden)
		return
	}

	eventItem, err := model.GetEventItemByRefId(ctx, z.Db, eventItemRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event item")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
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
	if htmx.Hx(r).Target() == "modalbody" {
		err = z.TemplateExecuteSub(w, "edit-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = z.TemplateExecute(w, "edit-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) CreateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		log.Debug().Err(err).Msg("bad event ref-id")
		z.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for event")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		z.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		z.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		log.Debug().Err(err).Msg("description form element empty")
		z.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	_, err = model.NewEventItem(ctx, z.Db, event.Id, description)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefId), http.StatusSeeOther)
}

func (z *ZHandler) UpdateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		z.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, eventRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		z.Error(w, "access denied", http.StatusForbidden)
		return
	}

	eventItem, err := model.GetEventItemByRefId(ctx, z.Db, eventItemRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
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
	earmark, err := model.GetEarmarkByEventItem(ctx, z.Db, eventItem)
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
	err = eventItem.Save(ctx, z.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", eventRefId)) {
		w.Header().Add("HX-Refresh", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (z *ZHandler) DeleteEventItem(w http.ResponseWriter, r *http.Request) {
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

	event, err := model.GetEventByRefId(ctx, z.Db, eventRefId)
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

	eventItem, err := model.GetEventItemByRefId(ctx, z.Db, eventItemRefId)
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

	err = eventItem.Delete(ctx, z.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
