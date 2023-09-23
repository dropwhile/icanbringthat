package zhandler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

func (z *ZHandler) ShowCreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"title":          "Create Account",
		"flashes":        z.SessMgr.FlashPopKey(ctx, "operations"),
		"next":           r.FormValue("next"),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = z.TemplateExecute(w, "create-account-form.gohtml", tplVars)
	if err != nil {
		log.Debug().Err(err).Msg("template error")
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) ShowSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"user":           user,
		"title":          "Settings",
		"flashes":        z.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = z.TemplateExecute(w, "show-settings.gohtml", tplVars)
	if err != nil {
		log.Debug().Err(err).Msg("template error")
		z.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (z *ZHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form data")
		z.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	passwd := r.PostFormValue("password")
	confirm_passwd := r.PostFormValue("confirm_password")
	if email == "" || name == "" || passwd == "" {
		z.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	if passwd != confirm_passwd {
		z.Error(w, "password and confirm_password fields do not match", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// see if user is logged in
	if auth.IsLoggedIn(ctx) {
		// got a user, you can't make a user if you are already logged
		// in!
		z.Error(w, "already logged in as a user", http.StatusForbidden)
		return
	}

	user, err := model.NewUser(ctx, z.Db, email, name, []byte(passwd))
	if err != nil {
		log.Error().Err(err).Msg("error adding user")
		z.Error(w, "error adding user", http.StatusBadRequest)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = z.SessMgr.RenewToken(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error renewing session token")
		z.Error(w, err.Error(), 500)
		return
	}
	// Then make the privilege-level change.
	z.SessMgr.Put(r.Context(), "user-id", user.Id)
	z.SessMgr.FlashAppend(ctx, "operations", "Account created. You are now logged in.")

	target := "/dashboard"
	if r.PostFormValue("next") != "" {
		target = r.FormValue("next")
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (z *ZHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form data")
		z.Error(w, err.Error(), http.StatusBadRequest)
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
		warnings = append(warnings, "Same Email specified was already present")
	}

	name := r.PostFormValue("name")
	if name != "" && name != user.Name {
		user.Name = name
		changes = true
		operations = append(operations, "Name update successfull")
	} else if name == user.Name {
		warnings = append(warnings, "Same Name specified was already present")
	}

	newPasswd := r.PostFormValue("password")
	if newPasswd != "" {
		confirmPassword := r.PostFormValue("confirm_password")
		if newPasswd != confirmPassword {
			warnings = append(warnings, "New Password and Confirm Password do not match")
		} else {
			oldPasswd := r.PostFormValue("old_password")
			if ok, err := user.CheckPass(ctx, []byte(oldPasswd)); err != nil || !ok {
				warnings = append(warnings, "Old Password invalid")
			} else {
				err = user.SetPass(ctx, []byte(newPasswd))
				if err != nil {
					log.Error().Err(err).Msg("error setting user password")
					z.Error(w, "error updating user", http.StatusInternalServerError)
					return
				}
				operations = append(operations, "Password update successfull")
				changes = true
			}
		}
	}

	if changes {
		err = user.Save(ctx, z.Db)
		if err != nil {
			log.Error().Err(err).Msg("error updating user")
			z.Error(w, "error updating user", http.StatusInternalServerError)
			return
		}
		z.SessMgr.FlashAppend(ctx, "operations", operations...)
	} else {
		z.SessMgr.FlashAppend(ctx, "errors", warnings...)
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (z *ZHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("bad session data")
		z.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	err = user.Delete(ctx, z.Db)
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		z.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	// destroy session
	err = z.SessMgr.Destroy(ctx)
	if err != nil {
		// not a fatal error (do not return 500 to user), since the user deleted
		// sucessfully already. just log the oddity
		log.Error().Err(err).Msg("error destroying session")
	}
	w.WriteHeader(200)
}
