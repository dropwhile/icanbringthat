package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func (x *Handler) ListEarmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, errx := x.service.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	earmarkCount, errx := x.service.GetEarmarksCount(ctx, user.ID)
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
	earmarks, _, errx := x.service.GetEarmarksPaginated(
		ctx, user.ID, 10, offset, archived,
	)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventItemIDs := util.ToListByFunc(earmarks, func(em *model.Earmark) int {
		return em.EventItemID
	})
	eventItems, errx := x.service.GetEventItemsByIDs(ctx, eventItemIDs)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventIDs := util.ToListByFunc(eventItems, func(e *model.EventItem) int {
		return e.EventID
	})
	events, errx := x.service.GetEventsByIDs(ctx, eventIDs)
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
		"flashes":        x.sessMgr.FlashPopAll(ctx),
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

	eventRefID, err := service.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	eventItemRefID, err := service.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	event, errx := x.service.GetEvent(ctx, eventRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	eventItem, errx := x.service.GetEventItem(ctx, eventItemRefID)
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

	eventRefID, err := service.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	eventItemRefID, err := service.ParseEventItemRefID(chi.URLParam(r, "iRefID"))
	if err != nil {
		x.BadRefIDError(w, "event-item", err)
		return
	}

	// check to ensure routing param exists for event
	event, errx := x.service.GetEvent(ctx, eventRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	// check to ensure routing param exists for event-item
	eventItem, errx := x.service.GetEventItem(ctx, eventItemRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.DBError(w, errx)
		}
		return
	}

	// ensure routing params are actually related/linked
	if eventItem.EventID != event.ID {
		slog.InfoContext(ctx, "eventItem.EventID and event.ID mismatch",
			slog.Int("user.ID", user.ID),
			slog.Int("event.ID", event.ID),
			slog.Int("eventItem.EventID", eventItem.EventID),
		)
		x.NotFoundError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	// ok for note to be empty
	note := r.FormValue("note")

	_, errx = x.service.NewEarmark(ctx, user, eventItem.ID, note)
	if errx != nil {
		switch errx.Code() {
		case errs.PermissionDenied:
			x.ForbiddenError(w, errx.Msg())
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

	refID, err := service.ParseEarmarkRefID(chi.URLParam(r, "mRefID"))
	if err != nil {
		x.BadRefIDError(w, "earmark", err)
		return
	}

	errx := x.service.DeleteEarmarkByRefID(ctx, user.ID, refID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			slog.InfoContext(ctx, "permission denied",
				slog.Int("user.ID", user.ID),
				logger.Err(errx),
			)
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
