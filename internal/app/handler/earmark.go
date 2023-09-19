package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, h.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
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
	earmarks, err := model.GetEarmarksByUserPaginated(ctx, h.Db, user, 10, offset)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	for i, em := range earmarks {
		ei, err := em.GetEventItem(ctx, h.Db)
		if err != nil {
			log.Info().Err(err).Msg("db error")
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		e, err := ei.GetEvent(ctx, h.Db)
		if err != nil {
			log.Info().Err(err).Msg("db error")
			http.Error(w, "db error", http.StatusInternalServerError)
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
	err = h.TemplateExecute(w, "list-earmarks.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowCreateEarmarkForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
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

	eventItem, err := model.GetEventItemByRefId(ctx, h.Db, eventItemRefId)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
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
	if Hx(r).Target() == "modalbody" {
		err = h.TemplateExecuteSub(w, "create-earmark-form.gohtml", "form", tplVars)
	} else {
		err = h.TemplateExecute(w, "create-earmark-form.gohtml", tplVars)
	}
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateEarmark(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "bad event-ref-id", http.StatusBadRequest)
		return
	}

	eventItemRefId, err := model.EventItemRefIdT.Parse(chi.URLParam(r, "iRefId"))
	if err != nil {
		log.Debug().Err(err).Msg("bad eventitem ref-id")
		http.Error(w, "bad eventitem-ref-id", http.StatusBadRequest)
		return
	}

	_, err = model.GetEventByRefId(ctx, h.Db, eventRefId)
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

	eventItem, err := model.GetEventItemByRefId(ctx, h.Db, eventItemRefId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for event_item")
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// make sure no earmark exists yet
	_, err = model.GetEarmarkByEventItem(ctx, h.Db, eventItem)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// good. this is what we want
	case err == nil:
		// earmark already exists!
		http.Error(w, "already earmarked by other user - access denied", http.StatusForbidden)
		return
	default:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ok for note to be empty
	note := r.FormValue("note")

	_, err = model.NewEarmark(ctx, h.Db, eventItem.Id, user.Id, note)
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

func (h *Handler) DeleteEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EarmarkRefIdT.Parse(chi.URLParam(r, "mRefId"))
	if err != nil {
		log.Debug().Err(err).Msg("bad earmark ref-id")
		http.Error(w, "bad earmark ref-id", http.StatusBadRequest)
		return
	}

	earmark, err := model.GetEarmarkByRefId(ctx, h.Db, refId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != earmark.UserId {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = earmark.Delete(ctx, h.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if Hx(r).CurrentUrl().HasPathPrefix("/events/") {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
