package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/somerr"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
)

func (x *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
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

	eventCount, err := model.GetEventCountsByUser(ctx, x.Db, user.ID)
	if err != nil {
		x.DBError(w, err)
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
	events, err := model.GetEventsByUserPaginatedFiltered(
		ctx, x.Db, user.ID, 10, offset*10, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		events = []*model.Event{}
	case err != nil:
		x.DBError(w, err)
		return
	}

	eventIDs := util.ToListByFunc(events, func(e *model.Event) int {
		return e.ID
	})
	eventItemCounts, err := model.GetEventItemsCountByEventIDs(ctx, x.Db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event items")
		eventItemCounts = []*model.EventItemCount{}
	case err != nil:
		x.DBError(w, err)
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
		"flashes":         x.SessMgr.FlashPopAll(ctx),
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
	if htmx.Hx(r).Target() == "eventCards" {
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

	notifCount, err := model.GetNotificationCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	owner := user.ID == event.UserID

	eventItems, err := model.GetEventItemsByEvent(ctx, x.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event items")
		eventItems = []*model.EventItem{}
	case err != nil:
		x.DBError(w, err)
		return
	}

	// sort if needed
	if len(event.ItemSortOrder) > 0 {
		log.Trace().
			Str("sortOrder", fmt.Sprintf("%v", event.ItemSortOrder)).
			Msg("sorting")
		sortSet := util.ToSetIndexed(event.ItemSortOrder)
		eventItemLen := len(eventItems)
		sortedList := make([]*model.EventItem, len(event.ItemSortOrder))
		unsortedList := make([]*model.EventItem, 0)
		for j := range eventItems {
			if idx, ok := sortSet[eventItems[j].ID]; ok && idx < eventItemLen {
				sortedList[idx] = eventItems[j]
			} else {
				unsortedList = append(unsortedList, eventItems[j])
			}
		}
		eventItems = append(unsortedList, sortedList...)
	}

	earmarks, err := model.GetEarmarksByEvent(ctx, x.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for earmarks")
		earmarks = []*model.Earmark{}
	case err != nil:
		x.DBError(w, err)
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
	earmarkUsers, err := model.GetUsersByIDs(ctx, x.Db, userIDs)
	if err != nil {
		x.DBError(w, err)
		return
	}
	has_favorite := false
	_, err = model.GetFavoriteByUserEvent(ctx, x.Db, user.ID, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		has_favorite = false
	case err != nil:
		x.DBError(w, err)
		return
	case err == nil:
		has_favorite = true
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
		"favorite":        has_favorite,
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
	if htmx.Hx(r).Target() == "modalbody" {
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	event, errx := service.GetEvent(ctx, x.Db, user.ID, refID)
	if errx != nil {
		switch errx.Code() {
		case somerr.NotFound:
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
	if htmx.Hx(r).Target() == "modalbody" {
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
		log.Debug().Err(err).Msg("error loading tz")
		tz = "Etc/UTC"
		loc, _ = time.LoadLocation(tz)
	}

	startTime, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing start time")
		x.BadFormDataError(w, err, "when", "loc")
		return
	}

	event, errx := service.CreateEvent(ctx, x.Db, user,
		name, description, startTime, tz)
	if errx != nil {
		switch errx.Code() {
		case somerr.InvalidArgument:
			x.BadFormDataError(w, err, errx.Meta("argument"))
		case somerr.PermissionDenied:
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	var (
		name        *string
		description *string
		tz          *string
		startTime   *time.Time
	)

	pname := r.PostFormValue("name")
	if pname != "" {
		name = &pname
	}
	pdesc := r.PostFormValue("description")
	if pdesc != "" {
		description = &pdesc
	}
	ptz := r.PostFormValue("timezone")
	when := r.PostFormValue("when")
	if when != "" && ptz != "" {
		loc, err := time.LoadLocation(ptz)
		if err != nil {
			log.Debug().Err(err).Msg("error loading tz")
			ptz = "Etc/UTC"
			loc, _ = time.LoadLocation(ptz)
		}
		t, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
		if err != nil {
			log.Debug().Err(err).Msg("error parsing start time")
			x.BadFormDataError(w, err, "when")
			return
		}
		lc := loc.String()
		tz = &lc
		startTime = &t
	}

	_, errx := service.UpdateEvent(ctx, x.Db, user.ID,
		refID, name, description, startTime, tz,
	)
	if errx != nil {
		switch errx.Code() {
		case somerr.NotFound:
			x.NotFoundError(w)
		case somerr.FailedPrecondition:
			x.BadFormDataError(w, errx)
		case somerr.PermissionDenied:
			x.AccessDeniedError(w)
		case somerr.InvalidArgument:
			x.BadFormDataError(w, errx, errx.Meta("argument"))
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%s", refID), http.StatusSeeOther)
}

// TODO:
func (x *Handler) UpdateEventItemSorting(w http.ResponseWriter, r *http.Request) {
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

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.AccessDeniedError(w)
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

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	sortOrder, ok := r.Form["sortOrder"]
	if !ok {
		x.BadFormDataError(w, err, "sortOrder")
		return
	}
	order := make([]int, 0)
	// make sure values are ok
	for _, v := range sortOrder {
		if i, err := strconv.Atoi(v); err != nil {
			x.BadFormDataError(w, err, "sortOrder")
			return
		} else {
			order = append(order, i)
		}
	}
	event.ItemSortOrder = util.Uniq(order)
	err = model.UpdateEvent(
		ctx, x.Db, event.ID,
		event.Name, event.Description,
		event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
	)
	if err != nil {
		x.DBError(w, err)
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.BadRefIDError(w, "event", err)
		return
	}

	if errx := service.DeleteEvent(ctx, x.Db, user.ID, refID); errx != nil {
		switch errx.Code() {
		case somerr.NotFound:
			x.NotFoundError(w)
		case somerr.PermissionDenied:
			x.AccessDeniedError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	if htmx.Hx(r).Request() {
		if htmx.Hx(r).CurrentUrl().HasPathPrefix(fmt.Sprintf("/events/%s", refID)) {
			x.SessMgr.FlashAppend(ctx, "success", "Event deleted.")
			w.Header().Add("HX-Redirect", "/events")
		} else {
			w.Header().Add("HX-Trigger-After-Swap", "count-updated")
		}
	}
	w.WriteHeader(http.StatusOK)
}
