package handler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ShowCreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to /account
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"title":          "Create Account",
		"flashes":        h.SessMgr.FlashPopKey(ctx, "operations"),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = h.TemplateExecute(w, "create-account-form.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	passwd := r.PostFormValue("password")
	confirm_passwd := r.PostFormValue("confirm_password")
	if email == "" || name == "" || passwd == "" {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	if passwd != confirm_passwd {
		http.Error(w, "password and confirm_password fields do not match", http.StatusBadRequest)
		return

	}

	ctx := r.Context()
	// see if user is logged in
	if auth.IsLoggedIn(ctx) {
		// got a user, you can't make a user if you are already logged
		// in!
		http.Error(w, "already logged in as a user", http.StatusForbidden)
		return
	}

	user, err := model.NewUser(ctx, h.Db, email, name, []byte(passwd))
	if err != nil {
		log.Info().Err(err).Msg("error adding user")
		http.Error(w, "error adding user", http.StatusBadRequest)
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
	h.SessMgr.FlashAppend(ctx, "operations", "Account created. You are now logged in.")

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) ShowSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"user":           user,
		"title":          "Settings",
		"flashes":        h.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = h.TemplateExecute(w, "show-settings.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	changes := false
	warnings := make([]string, 0)
	operations := make([]string, 0)

	email := r.PostFormValue("email")
	if email != "" && email != user.Email {
		user.Email = email
		changes = true
		operations = append(operations, "Email update successfull")
	} else if email == user.Email {
		warnings = append(warnings, "Same Email specified as was already present")
	}

	name := r.PostFormValue("name")
	if name != "" && name != user.Name {
		user.Name = name
		changes = true
		operations = append(operations, "Name update successfull")
	} else if name == user.Name {
		warnings = append(warnings, "Same Name specified as was already present")
	}

	newPasswd := r.PostFormValue("password")
	if newPasswd != "" {
		confirmPassword := r.PostFormValue("confirm_password")
		if newPasswd != confirmPassword {
			warnings = append(warnings, "New Password and Confirm Password do not match")
		} else {
			oldPasswd := r.PostFormValue("old_password")
			if ok, err := user.CheckPass(ctx, []byte(oldPasswd)); err != nil || !ok {
				warnings = append(warnings, "Current Password invalid")
			} else {
				user.SetPass(ctx, []byte(newPasswd))
				operations = append(operations, "Password update successfull")
				changes = true
			}
		}
	}

	if changes {
		err = user.Save(ctx, h.Db)
		if err != nil {
			http.Error(w, "error updating user", http.StatusInternalServerError)
			return
		}
		h.SessMgr.FlashAppend(ctx, "operations", operations...)
	} else {
		h.SessMgr.FlashAppend(ctx, "errors", warnings...)
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (h *Handler) ShowForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to /account
	if err == nil {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"user":           user,
		"flashes":        h.SessMgr.FlashPopKey(ctx, "operations"),
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = h.TemplateExecute(w, "forgot-password.gohtml", tplVars)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		http.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	err = user.Delete(ctx, h.Db)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}
