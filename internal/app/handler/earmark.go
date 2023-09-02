package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (h *Handler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(h.Db, ctx, user)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := uint(1)
	maxPageNum := ((earmarkCount / 10) + 1)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseUint(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(((earmarkCount / 10) + 1), uint(v))
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
	earmarks, err := model.GetEarmarksByUserPaginated(h.Db, ctx, user, 10, offset)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	for i, em := range earmarks {
		ei, err := em.GetEventItem(h.Db, ctx)
		if err != nil {
			mlog.Infox("db error", mlog.A("err", err))
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		e, err := ei.GetEvent(h.Db, ctx)
		if err != nil {
			mlog.Infox("db error", mlog.A("err", err))
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		ei.Event = e
		earmarks[i].EventItem = ei
	}

	tplVars := map[string]any{
		"user":           user,
		"earmarks":       earmarks,
		"earmarkCount":   earmarkCount,
		"pageNum":        pageNum,
		"pagePrev":       pagePrev,
		"pageNext":       pageNext,
		"title":          "My Earmarks",
		"nav":            "earmarks",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "list-earmarks.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EarmarkRefIdT.Parse(chi.URLParam(r, "refId"))
	if err != nil {
		http.Error(w, "bad earmark ref-id", http.StatusBadRequest)
		return
	}

	earmark, err := model.GetEarmarkByRefId(h.Db, ctx, refId)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != earmark.UserId {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = earmark.Delete(h.Db, ctx)
	if err != nil {
		mlog.Infox("db error", mlog.A("err", err))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
