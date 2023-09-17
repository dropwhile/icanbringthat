package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, h.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(((earmarkCount / 10) + 1), int(v))
				fmt.Println(pageNum)
			}
		}
	}

	offset := pageNum - 1
	earmarks, err := model.GetEarmarksByUserPaginated(ctx, h.Db, user, 10, offset)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	for i, em := range earmarks {
		ei, err := em.GetEventItem(ctx, h.Db)
		if err != nil {
			log.Info().Err(err).Msg("db error")
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		e, err := ei.GetEvent(ctx, h.Db)
		if err != nil {
			log.Info().Err(err).Msg("db error")
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
		"pgInput":        resources.NewPgInput(earmarkCount, 10, pageNum, "/earmarks"),
		"title":          "My Earmarks",
		"nav":            "earmarks",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
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

	refId, err := model.EarmarkRefIdT.Parse(chi.URLParam(r, "mRefId"))
	if err != nil {
		http.Error(w, "bad earmark ref-id", http.StatusBadRequest)
		return
	}

	earmark, err := model.GetEarmarkByRefId(ctx, h.Db, refId)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != earmark.UserId {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = earmark.Delete(ctx, h.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if Hx(r).CurrentUrl().HasPathPrefix("/events/") {
		w.Header().Add("HX-Refresh", "true")
	}
	w.WriteHeader(http.StatusOK)
}
