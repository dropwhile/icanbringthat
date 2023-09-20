package handler

import (
	"errors"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	events, err := model.GetEventsByUserPaginated(ctx, h.Db, user, 10, 0)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for events")
		events = []*model.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventCount, err := model.GetEventCountByUser(ctx, h.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, h.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"user":           user,
		"title":          "Dashboard",
		"nav":            "dashboard",
		"events":         events,
		"eventCount":     eventCount,
		"earmarkCount":   earmarkCount,
		"flashes":        h.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = h.TemplateExecute(w, "dashboard.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
