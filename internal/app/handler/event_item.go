package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/htmx"
)

func (x *Handler) ShowCreateEventItemForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.BadFormDataError(w, err)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
		return
	}

	tplVars := MapSA{
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
		err = x.TemplateExecuteSub(w, "create-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "create-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowEventItemEditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	eventItemRefID, err := model.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
		return
	}

	eventItem, err := model.GetEventItemByRefID(ctx, x.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	tplVars := MapSA{
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
		err = x.TemplateExecuteSub(w, "edit-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "edit-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) CreateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
		return
	}

	if event.Archived {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("event is archived")
		x.AccessDeniedError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		x.BadFormDataError(w, err, "description")
		return
	}

	_, err = model.NewEventItem(ctx, x.Db, event.ID, description)
	if err != nil {
		x.DBError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefID), http.StatusSeeOther)
}

func (x *Handler) UpdateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	eventItemRefID, err := model.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
		return
	}

	if event.Archived {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("event is archived")
		x.AccessDeniedError(w)
		return
	}

	eventItem, err := model.GetEventItemByRefID(ctx, x.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if eventItem.EventID != event.ID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.ID", event.ID).
			Int("eventItem.EventID", eventItem.EventID).
			Msg("eventItem.EventID and event.ID mismatch")
		x.NotFoundError(w)
		return
	}

	// check if earmark exists, and is marked by someone else
	// if so, disallow editing in that case.
	earmark, err := model.GetEarmarkByEventItem(ctx, x.Db, eventItem.ID)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		x.DBError(w, err)
		return
	case err == nil:
		if earmark.UserID != user.ID {
			log.Info().
				Int("user.ID", user.ID).
				Int("earmark.UserID", earmark.UserID).
				Msg("user id mismatch")
			x.ForbiddenError(w, "earmarked by other user")
			return
		}
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		x.BadFormDataError(w, err, "description")
		return
	}

	eventItem.Description = description
	err = model.UpdateEventItem(ctx, x.Db, eventItem.ID, eventItem.Description)
	if err != nil {
		x.DBError(w, err)
		return
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", eventRefID)) {
		w.Header().Add("HX-Refresh", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (x *Handler) DeleteEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	eventItemRefID, err := model.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
		return
	}

	if event.Archived {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("event is archived")
		x.AccessDeniedError(w)
		return
	}

	eventItem, err := model.GetEventItemByRefID(ctx, x.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if eventItem.EventID != event.ID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.ID", event.ID).
			Int("eventItem.EventID", eventItem.EventID).
			Msg("eventItem.EventID and event.ID mismatch")
		x.NotFoundError(w)
		return
	}

	err = model.DeleteEventItem(ctx, x.Db, eventItem.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
