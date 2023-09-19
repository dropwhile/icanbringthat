package handler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
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
		"flashes":        h.SessMgr.FlashPopKey(ctx, "login"),
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	// err := h.Tpl.ExecuteTemplate(w, "login-form.gohtml", tplVars)
	err = h.TemplateExecute(w, "login-form.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")

	if email == "" || passwd == "" {
		log.Debug().Msg("missing form data")
		http.Error(w, "Missing form data", http.StatusBadRequest)
		return
	}

	// find user...
	user, err := model.GetUserByEmail(ctx, h.Db, email)
	if err != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		h.SessMgr.FlashAppend(ctx, "login", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// validate credentials...
	ok, err := user.CheckPass(ctx, []byte(passwd))
	if err != nil || !ok {
		log.Debug().Err(err).Msg("invalid credentials: pass check fail")
		h.SessMgr.FlashAppend(ctx, "login", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = h.SessMgr.RenewToken(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Then make the privilege-level change.
	h.SessMgr.Put(r.Context(), "user-id", user.Id)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.SessMgr.Clear(r.Context()); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	if err := h.SessMgr.RenewToken(r.Context()); err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
