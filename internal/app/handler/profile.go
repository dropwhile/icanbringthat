package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
)

func (x *Handler) ShowProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	// parse user-id url param
	profileUserRefID, err := model.ParseUserRefID(chi.URLParam(r, "uRefID"))
	if err != nil {
		x.BadRefIDError(w, "user", err)
		return
	}

	selfView := false
	var profileUser *model.User
	if user.RefID == profileUserRefID {
		selfView = true
		profileUser = user
	} else {
		u, errx := service.GetUser(ctx, x.Db, profileUserRefID)
		if errx != nil {
			x.NotFoundError(w)
			return
		}
		profileUser = u
	}

	tplVars := MapSA{
		"user":        user,
		"profileUser": profileUser,
		"selfView":    selfView,
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-profile.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}
