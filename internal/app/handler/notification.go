package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/resources"
)

func (x *Handler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
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
	notifs, err := model.GetNotificationsByUserPaginated(ctx, x.Db, user.ID, 10, offset*10)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for favorite events")
		notifs = []*model.Notification{}
	case err != nil:
		x.DBError(w, err)
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

	notifRefID, err := model.ParseNotificationRefID(chi.URLParam(r, "nRefID"))
	if err != nil {
		x.BadRefIDError(w, "notification", err)
		return
	}

	notif, err := model.GetNotificationByRefID(ctx, x.Db, notifRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != notif.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("notif.UserID", notif.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
		return
	}

	err = model.DeleteNotification(ctx, x.Db, notif.ID)
	if err != nil {
		x.DBError(w, err)
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

	err = model.DeleteNotificationsByUser(ctx, x.Db, user.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
