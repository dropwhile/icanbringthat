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
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
)

func (x *Handler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, errx := service.GetNotificationsCount(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	earmarkCount, errx := service.GetEarmarksCount(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	extraQargs := url.Values{}
	maxCount := earmarkCount.Current
	archiveParam := r.FormValue("archive")
	archived := false
	if archiveParam == "1" {
		maxCount = earmarkCount.Archived
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
	earmarks, _, errx := service.GetEarmarksPaginated(
		ctx, x.Db, user.ID, 10, offset, archived,
	)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventItemIDs := util.ToListByFunc(earmarks, func(em *model.Earmark) int {
		return em.EventItemID
	})
	eventItems, errx := service.GetEventItemsByIDs(ctx, x.Db, eventItemIDs)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventIDs := util.ToListByFunc(eventItems, func(e *model.EventItem) int {
		return e.EventID
	})
	events, errx := service.GetEventsByIDs(ctx, x.Db, eventIDs)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventsMap := util.ToMapIndexedByFunc(events,
		func(v *model.Event) (int, *model.Event) { return v.ID, v })
	eventItemsMap := util.ToMapIndexedByFunc(eventItems,
		func(v *model.EventItem) (int, *model.EventItem) { return v.ID, v })

	title := "My Earmarks"
	if archived {
		title += " (Archived)"
	}
	tplVars := MapSA{
		"user":           user,
		"earmarks":       earmarks,
		"earmarkCount":   earmarkCount,
		"events":         eventsMap,
		"eventItems":     eventItemsMap,
		"notifCount":     notifCount,
		"title":          title,
		"nav":            "earmarks",
		"flashes":        x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
		"pgInput": resources.NewPgInput(
			maxCount, 10,
			pageNum, "/earmarks",
			extraQargs,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "earmarkCards" {
		err = x.TemplateExecuteSub(w, "list-earmarks.gohtml", "earmark_cards", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-earmarks.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowCreateEarmarkForm(w http.ResponseWriter, r *http.Request) {
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

	eventItemRefID, err := model.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	event, errx := service.GetEvent(ctx, x.Db, eventRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	eventItem, errx := service.GetEventItem(ctx, x.Db, eventItemRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	tplVars := MapSA{
		"user":           user,
		"event":          event,
		"eventItem":      eventItem,
		"title":          "Create Earmark",
		"nav":            "create-earmark",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "create-earmark-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "create-earmark-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) CreateEarmark(w http.ResponseWriter, r *http.Request) {
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

	eventItemRefID, err := model.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	event, errx := service.GetEvent(ctx, x.Db, eventRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	// non-owner must be verified before earmarking.
	// it is fine for owner to self-earmark though
	if !user.Verified && event.UserID != user.ID {
		x.ForbiddenError(w,
			"Account must be verified before earmarking is allowed.")
		return
	}

	if event.Archived {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("event is archived")
		x.AccessDeniedError(w)
		return
	}

	eventItem, errx := service.GetEventItem(ctx, x.Db, eventItemRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	if eventItem.EventID != event.ID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.ID", event.ID).
			Int("eventItem.EventID", eventItem.EventID).
			Msg("eventItem.EventID and event.ID mismatch")
		x.NotFoundError(w)
		return
	}

	// make sure no earmark exists yet
	_, errx = service.GetEarmarkByEventItemID(ctx, x.Db, eventItem.ID)
	if errx != nil {
		if errx.Code() != errs.NotFound {
			x.DBError(w, errx)
			return
		}
	} else {
		// earmark already exists!
		x.ForbiddenError(w, "already earmarked - access denied")
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	// ok for note to be empty
	note := r.FormValue("note")

	_, errx = service.NewEarmark(ctx, x.Db, eventItem.ID, user.ID, note)
	if errx != nil {
		switch errx.Code() {
		case errs.AlreadyExists:
			x.ForbiddenError(w, "already earmarked - access denied")
		default:
			x.DBError(w, errx)
		}
		return
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", eventRefID)) {
		w.Header().Add("HX-Refresh", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (x *Handler) DeleteEarmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := model.ParseEarmarkRefID(chi.URLParam(r, "mRefID"))
	if err != nil {
		x.BadRefIDError(w, "earmark", err)
		return
	}

	errx := service.DeleteEarmarkByRefID(ctx, x.Db, user.ID, refID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			log.Info().
				Int("user.ID", user.ID).
				Err(errx).
				Msg("permission denied")
			x.AccessDeniedError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Trigger-After-Swap", "count-updated")
	}
	w.WriteHeader(http.StatusOK)
}
