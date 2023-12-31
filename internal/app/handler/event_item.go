package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func (x *Handler) ShowCreateEventItemForm(w http.ResponseWriter, r *http.Request) {
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

	event, errx := x.service.GetEvent(ctx, eventRefID)
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
		slog.InfoContext(ctx,
			"user id mismatch",
			slog.Int("user.ID", user.ID),
			slog.Int("event.UserID", event.UserID),
		)
		x.AccessDeniedError(w)
		return
	}

	tplVars := MapSA{
		"user":           user,
		"event":          event,
		"title":          "Create Event Item",
		"nav":            "create-event-item",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "create-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "create-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowEventItemEditForm(w http.ResponseWriter, r *http.Request) {
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
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	if user.ID != event.UserID {
		slog.InfoContext(ctx, "user id mismatch",
			slog.Int("user.ID", user.ID),
			slog.Int("event.UserID", event.UserID),
		)
		x.AccessDeniedError(w)
		return
	}

	eventItem, errx := x.service.GetEventItem(ctx, eventItemRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	tplVars := MapSA{
		"user":           user,
		"event":          event,
		"eventItem":      eventItem,
		"title":          "Edit Event Item",
		"nav":            "edit-event-item",
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Hx(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "edit-eventitem-form.gohtml", "form", tplVars)
	} else {
		err = x.TemplateExecute(w, "edit-eventitem-form.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) CreateEventItem(w http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		x.BadFormDataError(w, err, "description")
		return
	}

	_, errx := x.service.AddEventItem(ctx, user.ID, eventRefID, description)
	if errx != nil {
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

	http.Redirect(w, r, fmt.Sprintf("/events/%s", eventRefID), http.StatusSeeOther)
}

func (x *Handler) UpdateEventItem(w http.ResponseWriter, r *http.Request) {
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

	// get event so we can ensure that the routing is valid..
	// er. /xxxx/yyyy where yyyy is actually an item for xxxx
	// and not someone just putting in /yolo/yyyy and getting
	// the expected result
	event, errx := x.service.GetEvent(ctx, eventRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		x.BadFormDataError(w, err, "description")
		return
	}

	_, errx = x.service.UpdateEventItem(
		ctx, user.ID, eventItemRefID, description,
		func(ei *model.EventItem) bool {
			return ei.EventID != event.ID
		},
	)
	if errx != nil {
		slog.DebugContext(ctx,
			"failed to update eventitem",
			"error", errx)
		switch errx.Code() {
		case errs.FailedPrecondition:
			slog.InfoContext(ctx,
				"eventItem.EventID and event.ID mismatch",
				slog.Int("user.ID", user.ID),
				slog.Int("event.ID", event.ID),
			)
			x.NotFoundError(w)
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		default:
			x.InternalServerError(w, errx.Msg())
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

func (x *Handler) DeleteEventItem(w http.ResponseWriter, r *http.Request) {
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

	// get event so we can ensure that the routing is valid..
	// er. /xxxx/yyyy where yyyy is actually an item for xxxx
	// and not someone just putting in /yolo/yyyy and getting
	// the expected result
	event, errx := x.service.GetEvent(ctx, eventRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			x.NotFoundError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	errx = x.service.RemoveEventItem(
		ctx, user.ID, eventItemRefID,
		func(ei *model.EventItem) bool {
			slog.DebugContext(ctx,
				"eventItem.EventID and event.ID comparison",
				slog.Int("user.ID", user.ID),
				slog.Int("event.ID", event.ID),
				slog.Int("eventItem.EventID", ei.EventID),
			)
			return ei.EventID != event.ID
		},
	)
	if errx != nil {
		slog.DebugContext(ctx, "failed to remove eventitem",
			"error", errx)
		switch errx.Code() {
		case errs.FailedPrecondition:
			x.NotFoundError(w)
		case errs.NotFound:
			x.NotFoundError(w)
		case errs.PermissionDenied:
			x.AccessDeniedError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
