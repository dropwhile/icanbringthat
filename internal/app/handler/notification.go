// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/htmx"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

func (x *Handler) NotificationsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, err := x.svc.GetNotificationsCount(ctx, user.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	pageNum := 1
	maxPageNum := resources.CalculateMaxPageNum(notifCount, 10)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(maxPageNum, int(v))
			}
		}
	}

	offset := pageNum - 1
	notifs, _, errx := x.svc.GetNotificationsPaginated(ctx, user.ID, 10, offset*10)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	tplVars := MapSA{
		"user":       user,
		"notifs":     notifs,
		"notifCount": notifCount,
		"title":      "Notifications",
		"nav":        "notifications",
		"pgInput": resources.NewPgInput(
			notifCount, 10, pageNum, "/notifications", nil,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "notifCount" {
		err = x.TemplateExecuteSub(w, "list-notifications.gohtml", "notif_count", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-notifications.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) NotificationDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := service.ParseNotificationRefID(r.PathValue("nRefID"))
	if err != nil {
		x.BadRefIDError(w, "notification", err)
		return
	}

	errx := x.svc.DeleteNotification(ctx, user.ID, refID)
	if errx != nil {
		slog.InfoContext(ctx, "error deleting notification", "error", errx)
		switch errx.Code() {
		case errs.Internal:
			x.InternalServerError(w, errx.Info())
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

func (x *Handler) NotificationsDeleteAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	errx := x.svc.DeleteAllNotifications(ctx, user.ID)
	if errx != nil {
		slog.InfoContext(ctx, "error deleting all notifications", "error", errx)
		switch errx.Code() {
		case errs.Internal:
			x.InternalServerError(w, errx.Info())
		case errs.Unauthenticated:
			x.BadSessionDataError(w)
		default:
			x.InternalServerError(w, "unexpected error")
		}
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).IsRequest() {
		htmx.Response(w).HxLocation(htmx.Request(r).CurrentUrl().Path)
	}
	w.WriteHeader(http.StatusOK)
}
