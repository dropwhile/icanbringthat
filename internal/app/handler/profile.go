// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"net/http"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
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
	profileUserRefID, err := service.ParseUserRefID(r.PathValue("uRefID"))
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
		u, errx := x.svc.GetUser(ctx, profileUserRefID)
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
