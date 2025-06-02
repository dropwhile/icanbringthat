// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"net/http"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func (x *Handler) DashboardShow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, errx := x.svc.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	eventCount, errx := x.svc.GetEventsCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	earmarkCount, errx := x.svc.GetEarmarksCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	favoriteCount, errx := x.svc.GetFavoriteEventsCount(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	events, _, errx := x.svc.GetEventsComingSoonPaginated(ctx, user.ID, 10, 0)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int {
		return e.ID
	})
	eventItemCounts, errx := x.svc.GetEventItemsCount(ctx, eventIDs)
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
		"flashes":         x.sessMgr.FlashPopAll(ctx),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "dashboard.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}
