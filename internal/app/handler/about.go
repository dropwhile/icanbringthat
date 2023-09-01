// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"net/http"

	"github.com/dropwhile/icanbringthat/internal/htmx"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

func (x *Handler) ShowAbout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	tplVars := MapSA{
		"user":  user,
		"title": "About",
		"nav":   "about",
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	if htmx.Request(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "about.gohtml", "about", tplVars)
	} else {
		err = x.TemplateExecute(w, "about.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}
