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
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
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

	eventCount, err := model.GetEventCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	maxPageNum := resources.CalculateMaxPageNum(eventCount, 10)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(maxPageNum, int(v))
			}
		}
	}

	offset := pageNum - 1
	events, err := model.GetEventsByUserPaginated(ctx, x.Db, user.ID, 10, offset*10)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		events = []*model.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	for i := range events {
		items, err := model.GetEventItemsByEvent(ctx, x.Db, events[i].ID)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			log.Info().Err(err).Msg("no rows for event items")
			items = []*model.EventItem{}
		case err != nil:
			log.Info().Err(err).Msg("db error")
			x.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		events[i].Items = items
	}

	tplVars := map[string]any{
		"user":           user,
		"events":         events,
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, refID)
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

	eventDict := util.StructToMap(event)
	owner := user.ID == event.UserID

	eventItems, err := model.GetEventItemsByEvent(ctx, x.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event items")
		eventItems = []*model.EventItem{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// sort if needed
	itemsList := make([]map[string]interface{}, 0)
	if len(event.ItemSortOrder) > 0 {
		log.Debug().Str("sortOrder", fmt.Sprintf("%v", event.ItemSortOrder)).Msg("sorting")
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
	for i := range eventItems {
		itemsList = append(itemsList, util.StructToMap(eventItems[i]))
	}

	earmarks, err := model.GetEarmarksByEvent(ctx, x.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for earmarks")
		earmarks = []*model.Earmark{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// associate earmarks and event items
	// and also collect the user ids associated with
	// earmarks
	earmarksMap := make(map[int]*model.Earmark)
	userIDsMap := make(map[int]struct{})
	for i := range earmarks {
		earmarksMap[earmarks[i].EventItemID] = earmarks[i]
		userIDsMap[earmarks[i].UserID] = struct{}{}
	}

	// now get the list of usrs ids and fetch the associated users
	userIDs := util.Keys(userIDsMap)
	earmarkUsers, err := model.GetUsersByIDs(ctx, x.Db, userIDs)
	if err != nil {
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// now associate the users with the earmarks
	earmarkUsersMap := make(map[int]*model.User)
	for i := range earmarkUsers {
		earmarkUsersMap[earmarkUsers[i].ID] = earmarkUsers[i]
	}

	for i := range eventItems {
		eid := eventItems[i].ID
		if earmark, ok := earmarksMap[eid]; ok {
			eaMap := util.StructToMap(earmark)
			if uu, ok := earmarkUsersMap[earmark.UserID]; ok {
				eaMap["User"] = util.StructToMap(uu)
			}
			itemsList[i]["Earmark"] = eaMap
		}
	}

	eventDict["Items"] = itemsList

	has_favorite := false
	_, err = model.GetFavoriteByUserEvent(ctx, x.Db, user.ID, event.ID)
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
		"event":          eventDict,
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, refID)
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

	event, err := model.NewEvent(ctx, x.Db, user.ID, name, description, startTime, &model.TimeZone{Location: loc})
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, refID)
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
		event.StartTimeTz = &model.TimeZone{Location: loc}
	}

	if name != "" {
		event.Name = name
	}
	if description != "" {
		event.Description = description
	}

	err = model.UpdateEvent(
		ctx, x.Db, event.ID,
		event.Name, event.Description,
		event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
	)
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

	eventRefID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event-ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
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
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
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
	order := make([]int, 0)
	// make sure values are ok
	for _, v := range sortOrder {
		if i, err := strconv.Atoi(v); err != nil {
			log.Debug().Err(err).Msg("bad form data")
			http.Error(w, "bad form data", http.StatusBadRequest)
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

	refID, err := model.ParseEventRefID(chi.URLParam(r, "eRefID"))
	if err != nil {
		x.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefID(ctx, x.Db, refID)
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
			Int("user.ID", user.ID).
			Int("event.UserID", event.UserID).
			Msg("user id mismatch")
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = model.DeleteEvent(ctx, x.Db, event.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
