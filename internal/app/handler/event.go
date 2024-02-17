package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func (x *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	notifCount, errx := x.svc.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventCount, errx := x.svc.GetEventsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	extraQargs := url.Values{}
	maxCount := eventCount.Current
	archiveParam := r.FormValue("archive")
	archived := false
	if archiveParam == "1" {
		maxCount = eventCount.Archived
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
	events, _, errx := x.svc.GetEventsPaginated(
		ctx, user.ID, 10, offset*10, archived)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int { return e.ID })
	eventItemCounts, errx := x.svc.GetEventItemsCount(ctx, eventIDs)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	eventItemCountsMap := util.ToMapIndexedByFunc(
		eventItemCounts,
		func(eic *model.EventItemCount) (int, int) {
			return eic.EventID, eic.Count
		})

	title := "My Events"
	if archived {
		title += " (Archived)"
	}
	tplVars := MapSA{
		"user":            user,
		"events":          events,
		"eventItemCounts": eventItemCountsMap,
		"eventCount":      eventCount,
		"notifCount":      notifCount,
		"title":           title,
		"nav":             "events",
		"flashes":         x.sessMgr.FlashPopAll(ctx),
		csrf.TemplateTag:  csrf.TemplateField(r),
		"csrfToken":       csrf.Token(r),
		"pgInput": resources.NewPgInput(
			maxCount, 10,
			pageNum, "/events",
			extraQargs,
		),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "eventCards" {
		err = x.TemplateExecuteSub(w, "list-events.gohtml", "event_cards", tplVars)
	} else {
		err = x.TemplateExecute(w, "list-events.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	notifCount, errx := x.svc.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	event, errx := x.svc.GetEvent(ctx, refID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	owner := user.ID == event.UserID

	eventItems, errx := x.svc.GetEventItemsByEventID(ctx, event.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	// sort if needed
	if len(event.ItemSortOrder) > 0 {
		// example of deferring a slow logging operation, and to avoid processing
		// if level is not debug
		// (note: this example isn't especially slow)
		slog.DebugContext(ctx, "item sorting",
			slog.Any("sortOrder",
				logger.DeferOperation(event.ItemSortOrder, func(i []int) string {
					return fmt.Sprintf("%v", i)
				})),
		)
		sortSet := util.ToSetIndexed(event.ItemSortOrder)
		sortedList := make([]*model.EventItem, len(event.ItemSortOrder))
		unsortedList := make([]*model.EventItem, 0)

		eventItemLen := len(eventItems)
		for j := range eventItems {
			if idx, ok := sortSet[eventItems[j].ID]; ok && idx < eventItemLen {
				sortedList[idx] = eventItems[j]
			} else {
				unsortedList = append(unsortedList, eventItems[j])
			}
		}
		// put unsorted (likely new) items at the front of the list
		eventItems = append(unsortedList, sortedList...)
	}

	earmarks, errx := x.svc.GetEarmarksByEventID(ctx, event.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	// associate earmarks and event items
	// and also collect the user ids associated with
	// earmarks
	userIDs := util.ToListByFunc(earmarks, func(e *model.Earmark) int {
		return e.UserID
	})
	userIDs = util.Uniq(userIDs)
	slices.Sort(userIDs)

	// now get the list of usrs ids and fetch the associated users
	earmarkUsers, errx := x.svc.GetUsersByIDs(ctx, userIDs)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	favorited := false
	_, errx = x.svc.GetFavoriteByUserEvent(ctx, user.ID, event.ID)
	if errx == nil {
		favorited = true
	} else {
		switch errx.Code() {
		case errs.NotFound:
			favorited = false
		case errs.Internal:
			x.DBError(w, errx)
			return
		}
	}

	earmarksMap := util.ToMapIndexedByFunc(earmarks,
		func(em *model.Earmark) (int, *model.Earmark) { return em.EventItemID, em },
	)

	earmarkUsersMap := util.ToMapIndexedByFunc(earmarkUsers,
		func(u *model.User) (int, *model.User) { return u.ID, u },
	)

	tplVars := MapSA{
		"user":            user,
		"owner":           owner,
		"event":           event,
		"eventItems":      eventItems,
		"earmarksMap":     earmarksMap,
		"earmarkUsersMap": earmarkUsersMap,
		"notifCount":      notifCount,
		"favorite":        favorited,
		"title":           "Event Details",
		"nav":             "show-event",
		csrf.TemplateTag:  csrf.TemplateField(r),
		"csrfToken":       csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-event.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowCreateEventForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	tplVars := MapSA{
		"user":           user,
		"title":          "Create Event",
		"nav":            "create-event",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "create-event-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "create-event-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowEditEventForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, errx := x.svc.GetEvent(ctx, refID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	if user.ID != event.UserID {
		x.AccessDeniedError(w)
		return
	}

	tplVars := MapSA{
		"user":           user,
		"event":          event,
		"title":          "Edit Event",
		"nav":            "edit-event",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "edit-event-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "edit-event-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" || description == "" || when == "" || tz == "" {
		x.BadFormDataError(w, err)
		return
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		slog.DebugContext(ctx, "error loading tz", "error", err)
		tz = "Etc/UTC"
		loc, _ = time.LoadLocation(tz)
	}

	startTime, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
	if err != nil {
		slog.DebugContext(ctx, "error parsing start time", "error", err)
		x.BadFormDataError(w, err, "when", "loc")
		return
	}

	event, errx := x.svc.CreateEvent(ctx, user, name, description, startTime, tz)
	if errx != nil {
		switch errx.Code() {
		case errs.InvalidArgument:
			x.BadFormDataError(w, err, errx.Meta("argument"))
		case errs.PermissionDenied:
			x.ForbiddenError(w, errx.Msg())
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefID), http.StatusSeeOther)
}

func (x *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	euvs := &service.EventUpdateValues{}

	name := r.PostFormValue("name")
	if name != "" {
		euvs.Name = mo.Some(name)
	}

	desc := r.PostFormValue("description")
	if desc != "" {
		euvs.Description = mo.Some(desc)
	}

	ptz := r.PostFormValue("timezone")
	when := r.PostFormValue("when")
	switch {
	case when == "" && ptz != "":
		x.BadFormDataError(w, nil, "when")
		return
	case when != "" && ptz == "":
		x.BadFormDataError(w, nil, "timezone")
		return
	case when != "" && ptz != "":
		loc, err := time.LoadLocation(ptz)
		if err != nil {
			x.BadFormDataError(w, nil, "timezone")
			return
		}
		t, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
		if err != nil {
			slog.DebugContext(ctx, "error parsing start time", "error", err)
			x.BadFormDataError(w, err, "when")
			return
		}
		euvs.StartTime = mo.Some(t)
		euvs.Tz = mo.Some(loc.String())
	}

	errx := x.svc.UpdateEvent(ctx, user.ID, refID, euvs)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.FailedPrecondition:
			x.BadFormDataError(w, errx)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		case errs.InvalidArgument:
			x.BadFormDataError(w, errx, errx.Meta("argument"))
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%s", refID), http.StatusSeeOther)
}

func (x *Handler) UpdateEventItemSorting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	eventRefID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	sortOrder, ok := r.Form["sortOrder"]
	if !ok {
		x.BadFormDataError(w, err, "sortOrder")
		return
	}

	// make sure values are ok
	order := make([]int, 0)
	for _, v := range sortOrder {
		if i, err := strconv.Atoi(v); err != nil {
			x.BadFormDataError(w, err, "sortOrder")
			return
		} else {
			order = append(order, i)
		}
	}
	order = util.Uniq(order)

	_, errx := x.svc.UpdateEventItemSorting(
		ctx, user.ID, eventRefID, order,
	)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		case errs.FailedPrecondition:
			x.BadFormDataError(w, errx, "sortOrder")
		case errs.Internal:
			x.DBError(w, errx)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (x *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	refID, err := service.ParseEventRefID(r.PathValue("eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	if errx := x.svc.DeleteEvent(ctx, user.ID, refID); errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	if htmx.Request(r).IsRequest() {
		if htmx.Request(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", refID)) {
			x.sessMgr.FlashAppend(ctx, "success", "Event deleted.")
			htmx.Response(w).HxLocation("/events")
		} else {
			htmx.Response(w).HxTriggerAfterSwap("count-updated")
		}
	}
	w.WriteHeader(http.StatusOK)
}
