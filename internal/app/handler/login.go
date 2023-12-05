package handler

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
)

func (x *Handler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tplVars := MapSA{
		"title":          "Login",
		"next":           r.FormValue("next"),
		"flashes":        x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "login-form.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")

	if email == "" || passwd == "" {
		x.BadFormDataError(w, nil, "missing form data")
		return
	}

	// find user...
	user, errx := service.GetUserByEmail(ctx, x.Db, email)
	if errx != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		x.SessMgr.FlashAppend(ctx, "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !user.PWAuth {
		// no valid auth flow
		log.Warn().
			Int("userID", user.ID).
			Msg("invalid credentials: no valid auth flow")
		x.SessMgr.FlashAppend(ctx, "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else {
		// validate credentials...
		ok, err := model.CheckPass(ctx, user.PWHash, []byte(passwd))
		if err != nil || !ok {
			log.Debug().Err(err).Msg("invalid credentials: pass check fail")
			x.SessMgr.FlashAppend(ctx, "error", "Invalid credentials")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = x.SessMgr.RenewToken(ctx)
	if err != nil {
		x.InternalServerError(w, "Session Error")
		return
	}
	// Then make the privilege-level change.
	x.SessMgr.Put(r.Context(), "user-id", user.ID)
	target := "/dashboard"
	if r.PostFormValue("next") != "" {
		target = r.FormValue("next")
	}
	x.SessMgr.FlashAppend(ctx, "success", "Login successful")
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (x *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := x.SessMgr.Clear(r.Context()); err != nil {
		x.InternalServerError(w, "Session Error")
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	if err := x.SessMgr.RenewToken(r.Context()); err != nil {
		x.InternalServerError(w, "Session Error")
		return
	}
	x.SessMgr.FlashAppend(ctx, "success", "Logout successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
