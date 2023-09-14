package handler

import (
	"net/http"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/gorilla/csrf"
)

func (h *Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	events, err := model.GetEventsByUserPaginated(ctx, h.Db, user, 10, 0)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventCount, err := model.GetEventCountByUser(ctx, h.Db, user)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, h.Db, user)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
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
		"flashes":        h.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "dashboard.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
