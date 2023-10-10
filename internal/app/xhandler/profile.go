package xhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/modelx"
)

func (x *XHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// parse user-id url param
	profileUserRefID, err := modelx.ParseUserRefID(chi.URLParam(r, "uRefID"))
	if err != nil {
		x.Error(w, "bad user ref-id", http.StatusNotFound)
		return
	}

	selfView := false
	var profileUser *modelx.User
	if user.RefID == profileUserRefID {
		selfView = true
		profileUser = user
	} else {
		profileUser, err = x.Query.GetUserByRefID(ctx, profileUserRefID)
		if err != nil {
			x.Error(w, "user not found", http.StatusNotFound)
			return
		}
	}

	tplVars := map[string]any{
		"user":        user,
		"profileUser": profileUser,
		"selfView":    selfView,
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-profile.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
