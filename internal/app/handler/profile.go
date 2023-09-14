package handler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// parse user-id url param
	profileUserRefId, err := model.UserRefIdT.Parse(chi.URLParam(r, "uRefId"))
	if err != nil {
		http.Error(w, "bad user ref-id", http.StatusBadRequest)
		return
	}

	selfView := false
	var profileUser *model.User
	if user.RefId == profileUserRefId {
		selfView = true
		profileUser = user
	} else {
		profileUser, err = model.GetUserByRefId(ctx, h.Db, profileUserRefId)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
	}

	tplVars := map[string]any{
		"user":        user,
		"profileUser": profileUser,
		"selfView":    selfView,
	}
	// render user profile view
	SetHeader("content-type", "text/html")
	err = h.TemplateExecute(w, "show-profile.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
