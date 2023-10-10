package xhandler

import (
	"errors"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/modelx"
)

func (x *XHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	events, err := x.Query.GetEventsComingSoonByUserPaginated(ctx, modelx.GetEventsComingSoonByUserPaginatedParams{
		UserID: user.ID,
		Limit:  10,
		Offset: 0,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for events")
		events = []*modelx.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventCount, err := x.Query.GetEventCountByUser(ctx, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	earmarkCount, err := x.Query.GetEarmarkCountByUser(ctx, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	favoriteCount, err := x.Query.GetFavoriteCountByUser(ctx, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
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
		"favoriteCount":  favoriteCount,
		"flashes":        x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "dashboard.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
