package xhandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/modelx"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/internal/util/htmx"
	"github.com/dropwhile/icbt/resources"
)

func (x *XHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventCount, err := x.Query.GetEventCountByUser(ctx, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := int64(1)
	maxPageNum := resources.CalculateMaxPageNum(eventCount, 10)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(maxPageNum, v)
			}
		}
	}

	offset := pageNum - 1
	events, err := x.Query.GetEventsByUserPaginated(ctx, modelx.GetEventsByUserPaginatedParams{
		UserID: user.ID,
		Limit:  10,
		Offset: int32(offset) * 10,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		events = []*modelx.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventsExpanded := make([]*modelx.EventExpanded, 0)
	for i := range events {
		items, err := x.Query.GetEventItemsByEvent(ctx, events[i].ID)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			log.Info().Err(err).Msg("no rows for event items")
			items = []*modelx.EventItem{}
		case err != nil:
			log.Info().Err(err).Msg("db error")
			x.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		expandedItems := make([]*modelx.EventItemExpanded, 0)
		for i := range items {
			expandedItems = append(expandedItems, &modelx.EventItemExpanded{
				EventItem: items[i],
			})
		}
		eventsExpanded = append(
			eventsExpanded,
			&modelx.EventExpanded{
				Event: events[i],
				Items: expandedItems,
			},
		)
	}

	tplVars := map[string]any{
		"user":           user,
		"events":         eventsExpanded,
		"eventCount":     eventCount,
		"pgInput":        resources.NewPgInput(eventCount, 10, pageNum, "/events"),
		"title":          "My Events",
		"nav":            "events",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}

	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "list-events.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) ShowEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	owner := user.ID == event.UserID

	eventItems, err := x.Query.GetEventItemsByEvent(ctx, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event items")
		eventItems = []*modelx.EventItem{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	eventExpanded := &modelx.EventExpanded{
		Event: event,
		Items: make([]*modelx.EventItemExpanded, 0),
	}
	eventItemsExpanded := make([]*modelx.EventItemExpanded, 0)
	// sort if needed
	if len(event.ItemSortOrder) > 0 {
		log.Debug().Str("sortOrder", fmt.Sprintf("%v", event.ItemSortOrder)).Msg("sorting")
		sortSet := util.ToSetIndexed(event.ItemSortOrder)
		eventItemLen := len(eventItems)
		sortedList := make([]*modelx.EventItemExpanded, len(event.ItemSortOrder))
		unsortedList := make([]*modelx.EventItemExpanded, 0)
		for j := range eventItems {
			if idx, ok := sortSet[eventItems[j].ID]; ok && idx < eventItemLen {
				sortedList[idx] = &modelx.EventItemExpanded{EventItem: eventItems[j], Event: nil}
			} else {
				unsortedList = append(unsortedList, &modelx.EventItemExpanded{EventItem: eventItems[j], Event: nil})
			}
		}
		eventItemsExpanded = append(unsortedList, sortedList...)
	}
	eventExpanded.Items = eventItemsExpanded

	earmarks, err := x.Query.GetEarmarksByEvent(ctx, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for earmarks")
		earmarks = []*modelx.Earmark{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// associate earmarks and event items
	// and also collect the user ids associated with
	// earmarks
	eventItemsMap := make(map[int32]*modelx.EventItemExpanded)
	for i := range eventItems {
		eventItemsMap[eventItems[i].ID] = &modelx.EventItemExpanded{EventItem: &modelx.EventItem{}}
	}

	earmarksExpanded := make([]*modelx.EarmarkExpanded, 0)
	userIdsMap := make(map[int32]struct{})
	for i := range earmarks {
		earmarksExpanded = append(earmarksExpanded, &modelx.EarmarkExpanded{Earmark: earmarks[i]})
		if ei, ok := eventItemsMap[earmarks[i].EventItemID]; ok {
			ei.Earmark = &modelx.EarmarkExpanded{Earmark: earmarks[i]}
			userIdsMap[earmarks[i].UserID] = struct{}{}
		}
	}

	// now get the list of usrs ids and fetch the associated users
	userIds := util.Keys(userIdsMap)
	earmarkUsers, err := x.Query.GetUsersByIds(ctx, userIds)
	if err != nil {
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// now associate the users with the earmarks
	earmarkUsersMap := make(map[int32]*modelx.User)
	for i := range earmarkUsers {
		earmarkUsersMap[earmarkUsers[i].ID] = earmarkUsers[i]
	}
	for i := range earmarks {
		if uu, ok := earmarkUsersMap[earmarks[i].UserID]; ok {
			earmarksExpanded[i].User = uu
		}
	}

	has_favorite := false
	_, err = x.Query.GetFavoriteByUserEvent(ctx, user.ID, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		has_favorite = false
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	case err == nil:
		has_favorite = true
	}

	tplVars := map[string]any{
		"user":           user,
		"owner":          owner,
		"event":          event,
		"favorite":       has_favorite,
		"title":          "Event Details",
		"nav":            "show-event",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-event.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) ShowCreateEventForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	tplVars := map[string]any{
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
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) ShowEditEventForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.ID != event.UserID {
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	tplVars := map[string]any{
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
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form data")
		x.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" || description == "" || when == "" || tz == "" {
		log.Debug().Msg("missing form data")
		x.Error(w, "bad form data", http.StatusBadRequest)
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
		x.Error(w, "bad form data - when", http.StatusBadRequest)
		return
	}

	event, err := x.Query.NewEvent(ctx, user.ID, name, description, startTime, loc.String())
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "error creating event", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefID), http.StatusSeeOther)
}

func (x *XHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.ID != event.UserID {
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" && description == "" && when == "" && tz == "" {
		x.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	switch {
	case when == "" && tz != "":
		fallthrough
	case when != "" && tz == "":
		x.Error(w, "bad form data", http.StatusBadRequest)
		return
	case when != "" && tz != "":
		loc, err := time.LoadLocation(tz)
		if err != nil {
			log.Debug().Err(err).Msg("error loading tz")
			tz = "Etc/UTC"
			loc, _ = time.LoadLocation(tz)
		}
		startTime, err := time.ParseInLocation("2006-01-02T15:04", when, loc)
		if err != nil {
			log.Debug().Err(err).Msg("error parsing start time")
			x.Error(w, "bad form data - when", http.StatusBadRequest)
			return
		}
		event.StartTime = startTime
		event.StartTimeTz = modelx.TimeZone{Location: loc}
	}

	if name != "" {
		event.Name = name
	}
	if description != "" {
		event.Description = description
	}

	updateParams := modelx.UpdateEventParams{
		Name:          &event.Name,
		Description:   &event.Description,
		ItemSortOrder: event.ItemSortOrder,
		StartTime: pgtype.Timestamptz{
			Time:  event.StartTime,
			Valid: true,
		},
		StartTimeTz: modelx.TimeZone{},
		ID:          event.ID,
	}
	err = x.Query.UpdateEvent(ctx, updateParams)
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "error updating event", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefID), http.StatusSeeOther)
}

func (x *XHandler) UpdateEventItemSorting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventRefID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int32("user.Id", user.ID).
			Int32("event.UserId", event.UserID).
			Msg("user id mismatch")
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sortOrder, ok := r.Form["sortOrder"]
	if !ok {
		log.Debug().Err(err).Msg("missing form data")
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	order := make([]int32, 0)
	// make sure values are ok
	for _, v := range sortOrder {
		if i, err := strconv.Atoi(v); err != nil {
			log.Debug().Err(err).Msg("bad form data")
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		} else {
			order = append(order, int32(i))
		}
	}
	updateParams := modelx.UpdateEventParams{
		ItemSortOrder: util.Uniq(order),
		ID:            event.ID,
	}
	err = x.Query.UpdateEvent(ctx, updateParams)
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "error updating event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (x *XHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refID, err := modelx.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := x.Query.GetEventByRefId(ctx, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.ID != event.UserID {
		log.Info().
			Int32("user.Id", user.ID).
			Int32("event.UserId", event.UserID).
			Msg("user id mismatch")
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = x.Query.DeleteEvent(ctx, event.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
