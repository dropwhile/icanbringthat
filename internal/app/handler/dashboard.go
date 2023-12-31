package handler

import (
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/middleware/auth"
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

	notifCount, errx := x.Service.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	eventCount, errx := x.Service.GetEventsCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	earmarkCount, errx := x.Service.GetEarmarksCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	favoriteCount, errx := x.Service.GetFavoriteEventsCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	events, _, errx := x.Service.GetEventsComingSoonPaginated(ctx, user.ID, 10, 0)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int {
		return e.ID
	})
	eventItemCounts, errx := x.Service.GetEventItemsCount(ctx, eventIDs)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
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
