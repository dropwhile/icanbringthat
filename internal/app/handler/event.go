package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventCount, err := model.GetEventCountByUser(ctx, h.Db, user)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	maxPageNum := ((eventCount / 10) + 1)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(maxPageNum, int(v))
				fmt.Println(pageNum)
			}
		}
	}
	pagePrev := 1
	if pageNum > 1 {
		pagePrev = pagePrev - 1
	}
	pageNext := maxPageNum
	if pageNum < maxPageNum {
		pageNext = pageNum + 1
	}

	offset := pageNum - 1
	events, err := model.GetEventsByUserPaginated(ctx, h.Db, user, 10, offset*10)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	for i := range events {
		items, err := model.GetEventItemsByEvent(ctx, h.Db, events[i])
		if err != nil {
			mlog.Infox("db error", mlog.A("err", err))
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		events[i].Items = items
	}

	tplVars := map[string]any{
		"user":           user,
		"events":         events,
		"eventCount":     eventCount,
		"pageNum":        pageNum,
		"pagePrev":       pagePrev,
		"pageNext":       pageNext,
		"title":          "My Events",
		"nav":            "events",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "list-events.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "refId"))
	if err != nil {
		http.Error(w, "bad event ref-id", http.StatusBadRequest)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, refId)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		fmt.Println(err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	owner := user.Id == event.UserId

	eventItems, err := model.GetEventItemsByEvent(ctx, h.Db, event)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
	}

	tplVars := map[string]any{
		"owner":      owner,
		"user":       user,
		"event":      event,
		"eventItems": eventItems,
	}
	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "event.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "refId"))
	if err != nil {
		http.Error(w, "bad event ref-id", http.StatusBadRequest)
		return
	}

	event, err := model.GetEventByRefId(ctx, h.Db, refId)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		mlog.Infof("user id mismatch %d != %d", user.Id, event.UserId)
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = event.Delete(ctx, h.Db)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
