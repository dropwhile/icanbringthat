package zhandler

import (
	"errors"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (z *ZHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	events, err := model.GetEventsComingSoonByUserPaginated(ctx, z.Db, user, 10, 0)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for events")
		events = []*model.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventCount, err := model.GetEventCountByUser(ctx, z.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, z.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
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
		"flashes":        z.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = z.TemplateExecute(w, "dashboard.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
