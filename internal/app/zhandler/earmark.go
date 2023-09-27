package zhandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util/htmx"
	"github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (z *ZHandler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, z.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(((earmarkCount / 10) + 1), int(v))
				fmt.Println(pageNum)
			}
		}
	}

	offset := pageNum - 1
	earmarks, err := model.GetEarmarksByUserPaginated(ctx, z.Db, user, 10, offset)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no earmarks")
		earmarks = []*model.Earmark{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	for i, em := range earmarks {
		ei, err := em.GetEventItem(ctx, z.Db)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			continue
		case err != nil:
			log.Info().Err(err).Msg("db error")
			z.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		e, err := ei.GetEvent(ctx, z.Db)
		// if no rows, or other db error, honk
		if err != nil {
			log.Info().Err(err).Msg("db error")
			z.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		ei.Event = e
		earmarks[i].EventItem = ei
	}

	tplVars := map[string]any{
		"user":           user,
		"earmarks":       earmarks,
		"earmarkCount":   earmarkCount,
		"pgInput":        resources.NewPgInput(earmarkCount, 10, pageNum, "/earmarks"),
		"title":          "My Earmarks",
		"nav":            "earmarks",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = z.TemplateExecute(w, "list-earmarks.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) ShowCreateEarmarkForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := model.EventRefIDT.Parse(chi.URLParam(r, "eRefID"))
	if err != nil {
		z.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefID, err := model.EventItemRefIDT.Parse(chi.URLParam(r, "iRefID"))
	if err != nil {
		z.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, z.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("event not found")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventItem, err := model.GetEventItemByRefID(ctx, z.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("event item not found")
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
		"title":          "Create Earmark",
		"nav":            "create-earmark",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "modalbody" {
		err = z.TemplateExecuteSub(w, "create-earmark-form.gohtml", "form", tplVars)
	} else {
		err = z.TemplateExecute(w, "create-earmark-form.gohtml", tplVars)
	}
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) CreateEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := model.EventRefIDT.Parse(chi.URLParam(r, "eRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad event ref-id")
		z.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	eventItemRefID, err := model.EventItemRefIDT.Parse(chi.URLParam(r, "iRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad eventitem ref-id")
		z.Error(w, "bad eventitem-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, z.Db, eventRefID)
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

	eventItem, err := model.GetEventItemByRefID(ctx, z.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for event_item")
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
		z.Error(w, "not found", http.StatusNotFound)
		return
	}

	// make sure no earmark exists yet
	_, err = model.GetEarmarkByEventItem(ctx, z.Db, eventItem)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// good. this is what we want
	case err == nil:
		// earmark already exists!
		z.Error(w, "already earmarked by other user - access denied", http.StatusForbidden)
		return
	default:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		z.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ok for note to be empty
	note := r.FormValue("note")

	_, err = model.NewEarmark(ctx, z.Db, eventItem.Id, user.Id, note)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", eventRefID)) {
		w.Header().Add("HX-Refresh", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (z *ZHandler) DeleteEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EarmarkRefIDT.Parse(chi.URLParam(r, "mRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad earmark ref-id")
		z.Error(w, "bad earmark ref-id", http.StatusNotFound)
		return
	}

	earmark, err := model.GetEarmarkByRefID(ctx, z.Db, refId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != earmark.UserId {
		z.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = earmark.Delete(ctx, z.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if htmx.Hx(r).CurrentUrl().HasPathPrefix("/events/") {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
