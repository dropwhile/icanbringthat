package zhandler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/go-chi/chi/v5"
)

func (z *ZHandler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// parse user-id url param
	profileUserRefId, err := model.UserRefIdT.Parse(chi.URLParam(r, "uRefId"))
	if err != nil {
		z.Error(w, "bad user ref-id", http.StatusNotFound)
		return
	}

	selfView := false
	var profileUser *model.User
	if user.RefId == profileUserRefId {
		selfView = true
		profileUser = user
	} else {
		profileUser, err = model.GetUserByRefId(ctx, z.Db, profileUserRefId)
		if err != nil {
			z.Error(w, "user not found", http.StatusNotFound)
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
	err = z.TemplateExecute(w, "show-profile.gohtml", tplVars)
	if err != nil {
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
