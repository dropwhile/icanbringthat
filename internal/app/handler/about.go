package handler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/htmx"
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
	if htmx.Hx(r).Target() == "modalbody" {
		err = x.TemplateExecuteSub(w, "about.gohtml", "about", tplVars)
	} else {
		err = x.TemplateExecute(w, "about.gohtml", tplVars)
	}
	if err != nil {
		x.TemplateError(w)
		return
	}
}