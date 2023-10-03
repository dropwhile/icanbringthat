package xhandler

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
)

func (x *XHandler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tplVars := map[string]any{
		"title":          "Login",
		"next":           r.FormValue("next"),
		"flashes":        x.SessMgr.FlashPopKey(ctx, "login"),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	// err := h.Tpl.ExecuteTemplate(w, "login-form.gohtml", tplVars)
	err = x.TemplateExecute(w, "login-form.gohtml", tplVars)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")

	if email == "" || passwd == "" {
		log.Debug().Msg("missing form data")
		x.Error(w, "Missing form data", http.StatusBadRequest)
		return
	}

	// find user...
	user, err := model.GetUserByEmail(ctx, x.Db, email)
	if err != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		x.SessMgr.FlashAppend(ctx, "login", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// validate credentials...
	ok, err := user.CheckPass(ctx, []byte(passwd))
	if err != nil || !ok {
		log.Debug().Err(err).Msg("invalid credentials: pass check fail")
		x.SessMgr.FlashAppend(ctx, "login", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = x.SessMgr.RenewToken(ctx)
	if err != nil {
		x.Error(w, err.Error(), 500)
		return
	}
	// Then make the privilege-level change.
	x.SessMgr.Put(r.Context(), "user-id", user.Id)
	target := "/dashboard"
	if r.PostFormValue("next") != "" {
		target = r.FormValue("next")
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (x *XHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := x.SessMgr.Clear(r.Context()); err != nil {
		x.Error(w, err.Error(), 500)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	if err := x.SessMgr.RenewToken(r.Context()); err != nil {
		x.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
