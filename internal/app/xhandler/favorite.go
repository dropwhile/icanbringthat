package xhandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util/htmx"
)

func (x *XHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := model.EventRefIDT.Parse(chi.URLParam(r, "eRefID"))
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
	if user.Id == event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id match")
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// check if already exists
	_, err = model.GetFavoriteByUserEvent(ctx, x.Db, user, event)
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

	_, err = model.NewFavorite(ctx, x.Db, user.Id, event.Id)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "favorite" {
		tplVars := map[string]any{
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

	eventRefID, err := model.EventRefIDT.Parse(chi.URLParam(r, "eRefID"))
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

	favorite, err := model.GetFavoriteByUserEvent(ctx, x.Db, user, event)
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

	err = favorite.Delete(ctx, x.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "favorite" {
		tplVars := map[string]any{
			"user":     user,
			"event":    event,
			"favorite": false,
		}
		if err := x.TemplateExecuteSub(w, "show-event.gohtml", "favorite", tplVars); err != nil {
			x.Error(w, "template error", http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
