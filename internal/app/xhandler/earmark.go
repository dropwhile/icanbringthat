package xhandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/modelx"
	"github.com/dropwhile/icbt/internal/util/htmx"
	"github.com/dropwhile/icbt/resources"
)

func (x *XHandler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	earmarkCount, err := x.Query.GetEarmarkCountByUser(ctx, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := int64(1)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(((earmarkCount / 10) + 1), v)
			}
		}
	}

	offset := pageNum - 1
	earmarks, err := x.Query.GetEarmarksByUserPaginated(ctx, modelx.GetEarmarksByUserPaginatedParams{
		UserID: user.ID,
		Limit:  10,
		Offset: int32(offset),
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no earmarks")
		earmarks = []*modelx.Earmark{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	earmarksExpanded := make([]*modelx.EarmarkExpanded, 0)
	for i, em := range earmarks {
		ei, err := x.Query.GetEventItemById(ctx, em.EventItemID)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			continue
		case err != nil:
			log.Info().Err(err).Msg("db error")
			x.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		e, err := x.Query.GetEventById(ctx, ei.EventID)
		// if no rows, or other db error, honk
		if err != nil {
			log.Info().Err(err).Msg("db error")
			x.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		eix := &modelx.EventItemExpanded{EventItem: ei, Event: e}
		emx := &modelx.EarmarkExpanded{Earmark: earmarks[i], EventItem: eix}
		earmarksExpanded = append(earmarksExpanded, emx)
	}

	tplVars := map[string]any{
		"user":           user,
		"earmarks":       earmarksExpanded,
		"earmarkCount":   earmarkCount,
		"pgInput":        resources.NewPgInput(earmarkCount, 10, pageNum, "/earmarks"),
		"title":          "My Earmarks",
		"nav":            "earmarks",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "list-earmarks.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) ShowCreateEarmarkForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefID, err := modelx.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("event not found")
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventItem, err := x.Query.GetEventItemByRefId(ctx, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("event item not found")
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	tplVars := map[string]any{
		"user":           user,
		"event":          event,
		"eventItem":      eventItem,
		"title":          "Create Earmark",
		"nav":            "create-earmark",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "create-earmark-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "create-earmark-form.gohtml", tplVars)
	}
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) CreateEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad event ref-id")
		x.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefID, err := modelx.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad eventitem ref-id")
		x.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for event")
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventItem, err := x.Query.GetEventItemByRefId(ctx, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for event_item")
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if eventItem.EventID != event.ID {
		log.Info().
			Int32("user.Id", user.ID).
			Int32("event.Id", event.ID).
			Int32("eventItem.EventId", eventItem.EventID).
			Msg("eventItem.EventId and event.Id mismatch")
		x.Error(w, "not found", http.StatusNotFound)
		return
	}

	// make sure no earmark exists yet
	_, err = x.Query.GetEarmarkByEventItem(ctx, eventItem.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// good. this is what we want
	case err == nil:
		// earmark already exists!
		x.Error(w, "already earmarked by other user - access denied", http.StatusForbidden)
		return
	default:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		x.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ok for note to be empty
	note := r.FormValue("note")

	_, err = x.Query.NewEarmark(ctx, eventItem.ID, user.ID, note)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", eventRefID)) {
		w.Header().Add("HX-Refresh", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (x *XHandler) DeleteEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refID, err := modelx.ParseEarmarkRefID(chi.URLParam(r, "mRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad earmark ref-id")
		x.Error(w, "bad earmark ref-id", http.StatusNotFound)
		return
	}

	earmark, err := x.Query.GetEarmarkByRefID(ctx, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.ID != earmark.UserID {
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = x.Query.DeleteEarmark(ctx, earmark.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if htmx.Hx(r).CurrentUrl().HasPathPrefix("/events/") {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
