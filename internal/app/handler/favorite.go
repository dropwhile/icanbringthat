// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/exp/slog"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/htmx"
	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func (x *Handler) FavoritesList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, errx := x.svc.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	favoriteCount, errx := x.svc.GetFavoriteEventsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
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
	events, _, errx := x.svc.GetFavoriteEventsPaginated(
		ctx, user.ID, 10, offset*10, archived,
	)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int {
		return e.ID
	})
	eventItemCounts, errx := x.svc.GetEventItemsCount(ctx, eventIDs)
	if errx != nil {
		x.DBError(w, errx)
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
		"flashes":         x.sessMgr.FlashPopAll(ctx),
		"pgInput": resources.NewPgInput(
			maxCount, 10,
			pageNum, "/favorites",
			extraQargs,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "favCards" {
		err = x.TemplateExecuteSub(w, "list-favorites.gohtml", "fav_cards", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-favorites.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) FavoriteAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, errx := x.svc.AddFavorite(ctx, user.ID, eventRefID)
	if errx != nil {
		slog.InfoContext(ctx, "error adding favorite", logger.Err(errx))
		switch errx.Code() {
		case errs.AlreadyExists:
			x.BadRequestError(w, "already favorited")
		case errs.Internal:
			x.InternalServerError(w, errx.Msg())
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		case errs.Unauthenticated:
			x.BadSessionDataError(w)
		default:
			x.InternalServerError(w, "unexpected error")
		}
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "favorite" {
		tplVars := MapSA{
			"user":     user,
			"event":    event,
			"favorite": true,
		}
		if err := x.TemplateExecuteSub(w, "show-event.gohtml", "favorite", tplVars); err != nil {
			x.TemplateError(w)
			return
		}
	} else {
		http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefID), http.StatusSeeOther)
	}
}

func (x *Handler) FavoriteDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	errx := x.svc.RemoveFavorite(ctx, user.ID, eventRefID)
	if errx != nil {
		slog.InfoContext(ctx, "error deleting favorite", logger.Err(errx))
		switch errx.Code() {
		case errs.Internal:
			x.InternalServerError(w, errx.Msg())
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		case errs.Unauthenticated:
			x.BadSessionDataError(w)
		default:
			x.InternalServerError(w, "unexpected error")
		}
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).IsRequest() {
		htmx.Response(w).HxTriggerAfterSwap("count-updated")
	}
	w.WriteHeader(http.StatusOK)
}
