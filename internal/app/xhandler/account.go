package xhandler

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
)

func (x *XHandler) ShowCreateAccount(w http.ResponseWriter, r *http.Request) {
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
		"flashes":        x.SessMgr.FlashPopAll(ctx),
		"next":           r.FormValue("next"),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "create-account-form.gohtml", tplVars)
	if err != nil {
		log.Debug().Err(err).Msg("template error")
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) ShowSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// parse user-id url param
	tplVars := map[string]any{
		"user":           user,
		"title":          "Settings",
		"flashes":        x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-settings.gohtml", tplVars)
	if err != nil {
		log.Debug().Err(err).Msg("template error")
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form data")
		x.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	passwd := r.PostFormValue("password")
	confirm_passwd := r.PostFormValue("confirm_password")
	if email == "" || name == "" || passwd == "" {
		x.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	if passwd != confirm_passwd {
		x.Error(w, "password and confirm_password fields do not match", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// see if user is logged in
	if auth.IsLoggedIn(ctx) {
		// got a user, you can't make a user if you are already logged
		// in!
		x.Error(w, "already logged in as a user", http.StatusForbidden)
		return
	}

	user, err := model.NewUser(ctx, x.Db, email, name, []byte(passwd))
	if err != nil {
		log.Error().Err(err).Msg("error adding user")
		x.Error(w, "error adding user", http.StatusBadRequest)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = x.SessMgr.RenewToken(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error renewing session token")
		x.Error(w, err.Error(), 500)
		return
	}
	// Then make the privilege-level change.
	x.SessMgr.Put(r.Context(), "user-id", user.ID)
	x.SessMgr.FlashAppend(ctx, "success", "Account created. You are now logged in.")

	target := "/dashboard"
	if r.PostFormValue("next") != "" {
		target = r.FormValue("next")
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (x *XHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Debug().Err(err).Msg("error parsing form data")
		x.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	changes := false
	warnings := make([]string, 0)
	successMsgs := make([]string, 0)

	email := r.PostFormValue("email")
	if email != "" && email != user.Email {
		user.Email = email
		user.Verified = false
		changes = true
		successMsgs = append(successMsgs, "Email update successfull")
	} else if email == user.Email {
		warnings = append(warnings, "Same Email specified was already present")
	}

	name := r.PostFormValue("name")
	if name != "" && name != user.Name {
		user.Name = name
		changes = true
		successMsgs = append(successMsgs, "Name update successfull")
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
			if ok, err := model.CheckPass(ctx, user, []byte(oldPasswd)); err != nil || !ok {
				warnings = append(warnings, "Old Password invalid")
			} else {
				pwhash, err := model.HashPass(ctx, []byte(newPasswd))
				if err != nil {
					log.Error().Err(err).Msg("error setting user password")
					x.Error(w, "error updating user", http.StatusInternalServerError)
					return
				}
				user.PWHash = pwhash
				successMsgs = append(successMsgs, "Password update successfull")
				changes = true
			}
		}
	}

	if changes {
		err = model.UpdateUser(ctx, x.Db, user.Email, user.Name, user.PWHash, user.Verified, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("error updating user")
			x.Error(w, "error updating user", http.StatusInternalServerError)
			return
		}
		x.SessMgr.FlashAppend(ctx, "success", successMsgs...)
	} else {
		x.SessMgr.FlashAppend(ctx, "error", warnings...)
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *XHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("bad session data")
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	err = model.DeleteUser(ctx, x.Db, user.ID)
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	// destroy session
	err = x.SessMgr.Destroy(ctx)
	if err != nil {
		// not a fatal error (do not return 500 to user), since the user deleted
		// sucessfully already. just log the oddity
		log.Error().Err(err).Msg("error destroying session")
	}
	w.WriteHeader(200)
}
