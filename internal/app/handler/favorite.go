package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/someerr"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
)

func (x *Handler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, errx := service.GetNotificationCount(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, err)
		return
	}

	favoriteCount, errx := service.GetFavoriteEventsCount(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, err)
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
	events, _, errx := service.GetFavoriteEventsPaginated(
		ctx, x.Db, user.ID, 10, offset*10, archived,
	)
	if errx != nil {
		x.DBError(w, err)
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int {
		return e.ID
	})
	eventItemCounts, errx := service.GetEventItemsCount(
		ctx, x.Db, user.ID, eventIDs)
	if errx != nil {
		x.DBError(w, err)
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
		"flashes":         x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag:  csrf.TemplateField(r),
		"csrfToken":       csrf.Token(r),
		"pgInput": resources.NewPgInput(
			maxCount, 10,
			pageNum, "/favorites",
			extraQargs,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "favCards" {
		err = x.TemplateExecuteSub(w, "list-favorites.gohtml", "fav_cards", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-favorites.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, errx := service.AddFavorite(ctx, x.Db, user.ID, eventRefID)
	if errx != nil {
		log.Info().
			Err(errx).
			Msg("error adding favorite")
		switch errx.Code() {
		case someerr.AlreadyExists:
			x.BadRequestError(w, "already favorited")
		case someerr.Internal:
			x.InternalServerError(w, errx.Msg())
		case someerr.NotFound:
			x.NotFoundError(w)
		case someerr.PermissionDenied:
			x.AccessDeniedError(w)
		case someerr.Unauthenticated:
			x.BadSessionDataError(w)
		default:
			x.InternalServerError(w, "unexpected error")
		}
		return
	}

	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "favorite" {
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

func (x *Handler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	errx := service.RemoveFavorite(ctx, x.Db, user.ID, eventRefID)
	if errx != nil {
		log.Info().
			Err(errx).
			Msg("error deleting favorite")
		switch errx.Code() {
		case someerr.Internal:
			x.InternalServerError(w, errx.Msg())
		case someerr.NotFound:
			x.NotFoundError(w)
		case someerr.PermissionDenied:
			x.AccessDeniedError(w)
		case someerr.Unauthenticated:
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
