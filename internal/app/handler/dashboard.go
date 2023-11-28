package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
)

func (x *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, err := model.GetNotificationCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	eventCount, err := model.GetEventCountsByUser(ctx, x.Db, user.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, x.Db, user)
	if err != nil {
		x.DBError(w, err)
		return
	}

	favoriteCount, err := model.GetFavoriteCountByUser(ctx, x.Db, user)
	if err != nil {
		x.DBError(w, err)
		return
	}

	events, err := model.GetEventsComingSoonByUserPaginated(ctx, x.Db, user.ID, 10, 0)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Msg("no rows for events")
		events = []*model.Event{}
	case err != nil:
		x.DBError(w, err)
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
		x.DBError(w, err)
		return
	}

	eventItemCountsMap := util.ToMapIndexedByFunc(
		eventItemCounts,
		func(eic *model.EventItemCount) (int, int) {
			return eic.EventID, eic.Count
		})

	// parse user-id url param
	tplVars := MapSA{
		"user":            user,
		"title":           "Dashboard",
		"nav":             "dashboard",
		"events":          events,
		"eventCount":      eventCount,
		"earmarkCount":    earmarkCount,
		"favoriteCount":   favoriteCount,
		"eventItemCounts": eventItemCountsMap,
		"notifCount":      notifCount,
		"flashes":         x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag:  csrf.TemplateField(r),
		"csrfToken":       csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "dashboard.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}
