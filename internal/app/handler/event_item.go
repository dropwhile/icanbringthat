package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ShowCreateEventItemForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	tplVars := map[string]any{
		"user":           user,
		"title":          "Create Event",
		"nav":            "create-event",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "create-event-form.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" || description == "" || when == "" || tz == "" {
		http.Error(w, "bad form data", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		tz = "Etc/UTC"
		loc, _ = time.LoadLocation(tz)
	}

	startTime, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
	if err != nil {
		http.Error(w, "bad form data - datetime", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	event, err := model.NewEvent(ctx, h.Db, user.Id, name, description, startTime, tz)
	if err != nil {
		http.Error(w, "error creating event", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefId), http.StatusSeeOther)
}

func (h *Handler) ShowEditEventItemForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event ref-id", http.StatusBadRequest)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, refId)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		fmt.Println(err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	tplVars := map[string]any{
		"user":           user,
		"event":          event,
		"title":          "Edit Event",
		"nav":            "edit-event",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "edit-event-form.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event ref-id", http.StatusBadRequest)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, refId)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		fmt.Println(err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" || description == "" || when == "" || tz == "" {
		http.Error(w, "bad form data", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		tz = "Etc/UTC"
		loc, _ = time.LoadLocation(tz)
	}

	startTime, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
	if err != nil {
		http.Error(w, "bad form data - datetime", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	event.Name = name
	event.Description = description
	event.StartTime = startTime
	event.StartTimeTZ = tz
	err = event.Save(ctx, h.Db)
	if err != nil {
		http.Error(w, "error updating event", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefId), http.StatusSeeOther)
}

func (h *Handler) DeleteEventItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusBadRequest)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		http.Error(w, "bad eventitem-ref-id", http.StatusBadRequest)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, eventRefId)
	if err != nil {
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
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
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
