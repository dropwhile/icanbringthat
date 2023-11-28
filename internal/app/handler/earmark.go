package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
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

	notifCount, err := model.GetNotificationCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	earmarkCount, err := model.GetEarmarkCountByUser(ctx, x.Db, user)
	if err != nil {
		x.DBError(w, err)
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
	earmarks, err := model.GetEarmarksByUserPaginatedFiltered(
		ctx, x.Db, user.ID, 10, offset, archived,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no earmarks")
		earmarks = []*model.Earmark{}
	case err != nil:
		x.DBError(w, err)
		return
	}

	eventItemIDs := util.ToListByFunc(earmarks, func(em *model.Earmark) int {
		return em.EventItemID
	})
	eventItems, err := model.GetEventItemsByIDs(ctx, x.Db, eventItemIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no event items")
		eventItems = []*model.EventItem{}
	case err != nil:
		x.DBError(w, err)
		return
	}

	eventIDs := util.ToListByFunc(eventItems, func(e *model.EventItem) int {
		return e.EventID
	})
	events, err := model.GetEventsByIDs(ctx, x.Db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no events")
		events = []*model.Event{}
	case err != nil:
		x.DBError(w, err)
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

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	eventItem, err := model.GetEventItemByRefID(ctx, x.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
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

	event, err := model.GetEventByRefID(ctx, x.Db, eventRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
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

	eventItem, err := model.GetEventItemByRefID(ctx, x.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
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
	_, err = model.GetEarmarkByEventItem(ctx, x.Db, eventItem.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// good. this is what we want
	case err == nil:
		// earmark already exists!
		x.ForbiddenError(w,
			"already earmarked by other user - access denied")
		return
	default:
		x.DBError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	// ok for note to be empty
	note := r.FormValue("note")

	_, err = model.NewEarmark(ctx, x.Db, eventItem.ID, user.ID, note)
	if err != nil {
		x.DBError(w, err)
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

	earmark, err := model.GetEarmarkByRefID(ctx, x.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
		return
	}

	if user.ID != earmark.UserID {
		x.AccessDeniedError(w)
		return
	}

	event, err := model.GetEventByEarmark(ctx, x.Db, earmark)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		x.NotFoundError(w)
		return
	case err != nil:
		x.DBError(w, err)
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

	err = model.DeleteEarmark(ctx, x.Db, earmark.ID)
	if err != nil {
		x.DBError(w, err)
		return
	}

	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Trigger-After-Swap", "count-updated")
	}
	w.WriteHeader(http.StatusOK)
}
