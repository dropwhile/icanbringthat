package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func (x *Handler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, err := service.GetNotificationsCount(ctx, x.Db, user.ID)
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
	notifs, _, errx := service.GetNotifcationsPaginated(ctx, x.Db, user.ID, 10, offset*10)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	tplVars := MapSA{
		"user":           user,
		"notifs":         notifs,
		"notifCount":     notifCount,
		"title":          "Notifications",
		"nav":            "notifications",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
		"pgInput": resources.NewPgInput(
			notifCount, 10, pageNum, "/notifications", nil,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "notifCount" {
		err = x.TemplateExecuteSub(w, "list-notifications.gohtml", "notif_count", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-notifications.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := model.ParseNotificationRefID(chi.URLParam(r, "nRefID"))
	if err != nil {
		x.BadRefIDError(w, "notification", err)
		return
	}

	errx := service.DeleteNotification(ctx, x.Db, user.ID, refID)
	if errx != nil {
		slog.InfoContext(ctx, "error deleting notification", "error", errx)
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
	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Trigger-After-Swap", "count-updated")
	}
	w.WriteHeader(http.StatusOK)
}

func (x *Handler) DeleteAllNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	errx := service.DeleteAllNotifications(ctx, x.Db, user.ID)
	if errx != nil {
		slog.InfoContext(ctx, "error deleting all notifications", "error", errx)
		switch errx.Code() {
		case errs.Internal:
			x.InternalServerError(w, errx.Msg())
		case errs.Unauthenticated:
			x.BadSessionDataError(w)
		default:
			x.InternalServerError(w, "unexpected error")
		}
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
