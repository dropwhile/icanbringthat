package zhandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (z *ZHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	eventCount, err := model.GetEventCountByUser(ctx, z.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	maxPageNum := resources.CalculateMaxPageNum(eventCount, 10)
	pageNumParam := r.FormValue("page")
	if pageNumParam != "" {
		if v, err := strconv.ParseInt(pageNumParam, 10, 0); err == nil {
			if v > 1 {
				pageNum = min(maxPageNum, int(v))
				fmt.Println(pageNum)
			}
		}
	}

	offset := pageNum - 1
	events, err := model.GetEventsByUserPaginated(ctx, z.Db, user, 10, offset*10)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		events = []*model.Event{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	for i := range events {
		items, err := model.GetEventItemsByEvent(ctx, z.Db, events[i])
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			log.Info().Err(err).Msg("no rows for event items")
			items = []*model.EventItem{}
		case err != nil:
			log.Info().Err(err).Msg("db error")
			z.Error(w, "db error", http.StatusInternalServerError)
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
	err = z.TemplateExecute(w, "list-events.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) ShowEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, refId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	owner := user.Id == event.UserId

	eventItems, err := model.GetEventItemsByEvent(ctx, z.Db, event)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event items")
		eventItems = []*model.EventItem{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	event.Items = eventItems

	earmarks, err := model.GetEarmarksByEvent(ctx, z.Db, event)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for earmarks")
		earmarks = []*model.Earmark{}
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// associate earmarks and event items
	// and also collect the user ids associated with
	// earmarks
	eventItemsMap := make(map[int]*model.EventItem)
	for i := range eventItems {
		eventItemsMap[eventItems[i].Id] = eventItems[i]
	}

	userIdsMap := make(map[int]struct{})
	for i := range earmarks {
		if ei, ok := eventItemsMap[earmarks[i].EventItemId]; ok {
			ei.Earmark = earmarks[i]
			userIdsMap[earmarks[i].UserId] = struct{}{}
		}
	}

	// now get the list of usrs ids and fetch the associated users
	userIds := util.Keys(userIdsMap)
	earmarkUsers, err := model.GetUsersByIds(ctx, z.Db, userIds)
	if err != nil {
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// now associate the users with the earmarks
	earmarkUsersMap := make(map[int]*model.User)
	for i := range earmarkUsers {
		earmarkUsersMap[earmarkUsers[i].Id] = earmarkUsers[i]
	}
	for i := range earmarks {
		if uu, ok := earmarkUsersMap[earmarks[i].UserId]; ok {
			earmarks[i].User = uu
		}
	}

	tplVars := map[string]any{
		"user":           user,
		"owner":          owner,
		"event":          event,
		"title":          "Event Details",
		"nav":            "show-event",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = z.TemplateExecute(w, "show-event.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) ShowCreateEventForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
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
	err = z.TemplateExecute(w, "create-event-form.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) ShowEditEventForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, refId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event")
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		z.Error(w, "access denied", http.StatusForbidden)
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
	err = z.TemplateExecute(w, "edit-event-form.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form data")
		z.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" || description == "" || when == "" || tz == "" {
		log.Debug().Msg("missing form data")
		z.Error(w, "bad form data", http.StatusBadRequest)
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
		z.Error(w, "bad form data - when", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	event, err := model.NewEvent(ctx, z.Db, user.Id, name, description, startTime, loc.String())
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		z.Error(w, "error creating event", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefId), http.StatusSeeOther)
}

func (z *ZHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, refId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Debug().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		z.Error(w, "access denied", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		z.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	when := r.PostFormValue("when")
	tz := r.PostFormValue("timezone")
	if name == "" && description == "" && when == "" && tz == "" {
		z.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	switch {
	case when == "" && tz != "":
		fallthrough
	case when != "" && tz == "":
		z.Error(w, "bad form data", http.StatusBadRequest)
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
			z.Error(w, "bad form data - when", http.StatusBadRequest)
			return
		}
		event.StartTime = startTime
		event.StartTimeTZ = loc.String()
	}

	if name != "" {
		event.Name = name
	}
	if description != "" {
		event.Description = description
	}

	err = event.Save(ctx, z.Db)
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		z.Error(w, "error updating event", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/events/%s", event.RefId), http.StatusSeeOther)
}

func (z *ZHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	refId, err := model.EventRefIdT.Parse(chi.URLParam(r, "eRefId"))
	if err != nil {
		z.Error(w, "bad event ref-id", http.StatusNotFound)
		return
	}

	event, err := model.GetEventByRefId(ctx, z.Db, refId)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		z.Error(w, "not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if user.Id != event.UserId {
		log.Info().
			Int("user.Id", user.Id).
			Int("event.UserId", event.UserId).
			Msg("user id mismatch")
		z.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = event.Delete(ctx, z.Db)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}