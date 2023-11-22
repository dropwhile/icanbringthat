package xhandler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
)

func (x *XHandler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	notifCount, err := model.GetNotificationCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	favoriteCount, err := model.GetFavoriteCountByUser(ctx, x.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	extraQargs := url.Values{}
	maxCount := favoriteCount.Current
	archiveParam := r.FormValue("archive")
	archived := false
	if archiveParam == "1" {
		maxCount = favoriteCount.Archived
		extraQargs.Add("archive", "1")
		archived = true
	}

	pageNum := 1
	maxPageNum := resources.CalculateMaxPageNum(maxCount, 10)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(maxPageNum, int(v))
			}
		}
	}

	offset := pageNum - 1
	events, err := model.GetFavoriteEventsByUserPaginatedFiltered(
		ctx, x.Db, user.ID, 10, offset*10, archived,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for favorite events")
		events = []*model.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int {
		return e.ID
	})
	eventItemCounts, err := model.GetEventItemsCountByEventIDs(ctx, x.Db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event items")
		eventItemCounts = []*model.EventItemCount{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventItemCountsMap := util.ToMapIndexedByFunc(
		eventItemCounts,
		func(eic *model.EventItemCount) (int, int) {
			return eic.EventID, eic.Count
		})

	title := "My Favorites"
	if archived {
		title += " (Archived)"
	}
	tplVars := MapSA{
		"user":            user,
		"events":          events,
		"favoriteCount":   favoriteCount,
		"eventItemCounts": eventItemCountsMap,
		"notifCount":      notifCount,
		"title":           title,
		"nav":             "favorites",
		"flashes":         x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag:  csrf.TemplateField(r),
		"csrfToken":       csrf.Token(r),
		"pgInput": resources.NewPgInput(
			maxCount, 10,
			pageNum, "/favorites",
			extraQargs,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "favCards" {
		err = x.TemplateExecuteSub(w, "list-favorites.gohtml", "fav_cards", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-favorites.gohtml", tplVars)
	}
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		log.Debug().Err(err).Msg("bad event ref-id")
		x.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
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

	// can't favorite your own event
	if user.ID == event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id match")
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// check if already exists
	_, err = model.GetFavoriteByUserEvent(ctx, x.Db, user.ID, event.ID)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	case err == nil:
		log.Info().Msg("already exists")
		x.Error(w, "already favorited", http.StatusBadRequest)
		return
	}

	_, err = model.CreateFavorite(ctx, x.Db, user.ID, event.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "favorite" {
		tplVars := MapSA{
			"user":     user,
			"event":    event,
			"favorite": true,
		}
		if err := x.TemplateExecuteSub(w, "show-event.gohtml", "favorite", tplVars); err != nil {
			x.Error(w, "template error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefID), http.StatusSeeOther)
	}
}

func (x *XHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	favorite, err := model.GetFavoriteByUserEvent(ctx, x.Db, user.ID, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Msg("favorite not found")
		x.Error(w, "not favorited", http.StatusBadRequest)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	err = model.DeleteFavorite(ctx, x.Db, favorite.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Trigger-After-Swap", "count-updated")
	}
	w.WriteHeader(http.StatusOK)
}
