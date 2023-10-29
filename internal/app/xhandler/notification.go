package xhandler

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
	"github.com/dropwhile/icbt/internal/util/htmx"
	"github.com/dropwhile/icbt/resources"
)

func (x *XHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	notifCount, err := model.GetNotificationCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
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
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	tplVars := map[string]any{
		"user":           user,
		"notifs":         notifs,
		"notifCount":     notifCount,
		"pgInput":        resources.NewPgInput(notifCount, 10, pageNum, "/notifications"),
		"title":          "Notifications",
		"nav":            "notifications",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "list-notifications.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	notifRefID, err := model.ParseNotificationRefID(chi.URLParam(r, "nRefID"))
	if err != nil {
		http.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	notif, err := model.GetNotificationByRefID(ctx, x.Db, notifRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.ID != notif.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("notif.UserID", notif.UserID).
			Msg("user id mismatch")
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = model.DeleteNotification(ctx, x.Db, notif.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func (x *XHandler) DeleteAllNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if notifCount, err := model.GetNotificationCountByUser(ctx, x.Db, user.ID); err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	} else if notifCount > 0 {
		err = model.DeleteNotificationsByUser(ctx, x.Db, user.ID)
		if err != nil {
			log.Info().Err(err).Msg("db error")
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).CurrentUrl().HasPathPrefix("/notifications") {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
