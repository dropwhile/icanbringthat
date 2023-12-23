package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func (x *Handler) ShowCreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// parse user-id url param
	tplVars := MapSA{
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
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ShowSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	credentials, errx := service.GetUserCredentialsByUser(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	apikey, errx := service.GetApiKeyByUser(ctx, x.Db, user.ID)
	if errx != nil {
		if errx.Code() != errs.NotFound {
			x.DBError(w, errx)
			return
		}
	}

	notifCount, errx := service.GetNotificationsCount(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	// parse user-id url param
	tplVars := MapSA{
		"user":           user,
		"credentials":    credentials,
		"apikey":         apikey,
		"title":          "Settings",
		"notifCount":     notifCount,
		"flashes":        x.SessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-settings.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	passwd := r.PostFormValue("password")
	confirm_passwd := r.PostFormValue("confirm_password")
	if email == "" || name == "" || passwd == "" {
		x.BadFormDataError(w, nil)
		return
	}
	if passwd != confirm_passwd {
		x.BadRequestError(w, "password and confirm_password fields do not match")
		return
	}

	ctx := r.Context()
	// see if user is logged in
	if auth.IsLoggedIn(ctx) {
		// got a user, you can't make a user if you are already logged
		// in!
		x.ForbiddenError(w, "already logged in as a user")
		return
	}

	user, errx := service.NewUser(ctx, x.Db, email, name, []byte(passwd))
	if errx != nil {
		slog.ErrorContext(ctx, "error adding user", logger.Err(errx))
		x.BadRequestError(w, "error adding user")
		return
	}

	_, errx = service.NewNotification(ctx, x.Db, user.ID,
		`Account is not currently verified. Please verify account in link:/settings.`,
	)
	if errx != nil {
		// this is a nonfatal error
		slog.ErrorContext(ctx, "error adding account notification",
			logger.Err(errx))
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err := x.SessMgr.RenewToken(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error renewing session token",
			logger.Err(err))
		x.InternalServerError(w, "Session Error")
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

func (x *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	changes := false
	warnings := make([]string, 0)
	successMsgs := make([]string, 0)
	updateVals := service.UserUpdateValues{}

	email := r.PostFormValue("email")
	if email != "" && email != user.Email {
		user.Email = email
		user.Verified = false
		changes = true
		updateVals.Email = mo.Some(email)
		updateVals.Verified = mo.Some(false)
		successMsgs = append(successMsgs, "Email update successfull")
	} else if email == user.Email {
		warnings = append(warnings, "Same Email specified was already present")
	}

	name := r.PostFormValue("name")
	if name != "" && name != user.Name {
		user.Name = name
		changes = true
		updateVals.Name = mo.Some(name)
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
			if ok, err := model.CheckPass(ctx, user.PWHash, []byte(oldPasswd)); err != nil || !ok {
				warnings = append(warnings, "Old Password invalid")
			} else {
				pwhash, err := model.HashPass(ctx, []byte(newPasswd))
				if err != nil {
					slog.ErrorContext(ctx, "error setting user password",
						logger.Err(err))
					x.InternalServerError(w, "error updating user")
					return
				}
				user.PWHash = pwhash
				updateVals.PWHash = mo.Some(pwhash)
				successMsgs = append(successMsgs, "Password update successfull")
				changes = true
			}
		}
	}

	if changes {
		errx := service.UpdateUser(ctx, x.Db, user.ID, updateVals)
		if errx != nil {
			slog.ErrorContext(ctx, "error updating user",
				logger.Err(errx))
			x.InternalServerError(w, "error updating user")
			return
		}
		x.SessMgr.FlashAppend(ctx, "success", successMsgs...)
	} else {
		x.SessMgr.FlashAppend(ctx, "error", warnings...)
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) UpdateAuthSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	changes := false
	authPW := r.PostFormValue("auth_passauth")
	authPK := r.PostFormValue("auth_passkeys")

	if authPW == "" && authPK == "" {
		x.BadFormDataError(w, nil, "auth params")
		return
	}

	// ensure we have at least one passkey first
	pkCount, errx := service.GetUserCredentialCountByUser(ctx, x.Db, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}
	hasPasskeys := false
	if pkCount > 0 {
		hasPasskeys = true
	}

	updateVals := service.UserUpdateValues{}
	switch authPW {
	case "off":
		if user.PWAuth {
			if !user.WebAuthn {
				x.SessMgr.FlashAppend(ctx, "error", "Account must be verified before enabling api access")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			user.PWAuth = false
			updateVals.PWAuth = mo.Some(false)
		}
	case "on":
		if !user.PWAuth {
			changes = true
			user.PWAuth = true
			updateVals.PWAuth = mo.Some(true)
		}
	case "":
		// nothing
	default:
		x.BadFormDataError(w, nil, "auth_passauth")
		return
	}

	switch authPK {
	case "off":
		if user.WebAuthn {
			if !user.PWAuth {
				x.SessMgr.FlashAppend(ctx, "error", "Refusing to disable passkey auth without alternative auth enabled")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			user.WebAuthn = false
			updateVals.WebAuthn = mo.Some(false)
		}
	case "on":
		if !user.WebAuthn {
			if !hasPasskeys {
				x.SessMgr.FlashAppend(ctx, "error",
					"Must have at least one passkey registered before enabling passkey auth")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			user.WebAuthn = true
			updateVals.WebAuthn = mo.Some(true)
		}
	case "":
		// nothing
	default:
		x.BadFormDataError(w, nil, "auth_passkeys")
		return
	}

	if !changes {
		x.SessMgr.FlashAppend(ctx, "error", "no changes made")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	errx = service.UpdateUser(ctx, x.Db, user.ID, updateVals)
	if errx != nil {
		slog.ErrorContext(ctx, "error updating user auth",
			logger.Err(errx))
		x.InternalServerError(w, "error updating user auth")
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) UpdateApiAuthSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	changes := false
	apiAccess := r.PostFormValue("api_access")
	rotateApiKey := r.PostFormValue("rotate_apikey")

	if apiAccess == "" && rotateApiKey == "" {
		x.BadFormDataError(w, nil, "bad params")
		return
	}

	updateVals := service.UserUpdateValues{}
	switch apiAccess {
	case "off":
		if user.ApiAccess {
			changes = true
			user.ApiAccess = false
			updateVals.ApiAccess = mo.Some(false)
		}
	case "on":
		if !user.ApiAccess {
			if !user.Verified {
				x.SessMgr.FlashAppend(ctx, "error", "Refusing to enable api access without a verified account")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			user.ApiAccess = true
			updateVals.ApiAccess = mo.Some(true)
			_, errx := service.NewApiKeyIfNotExists(ctx, x.Db, user.ID)
			if errx != nil {
				x.DBError(w, errx)
				return
			}
		}
	case "":
		// nothing
	default:
		x.BadFormDataError(w, nil, "api_access")
		return
	}

	if rotateApiKey == "true" {
		changes = true
	}

	if !changes {
		x.SessMgr.FlashAppend(ctx, "error", "no changes made")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	if rotateApiKey == "true" {
		if _, errx := service.NewApiKey(ctx, x.Db, user.ID); errx != nil {
			slog.ErrorContext(ctx, "error rotating api key",
				logger.Err(errx))
			x.InternalServerError(w, "error rotating api key")
			return
		}
	}

	if apiAccess != "" {
		if errx := service.UpdateUser(ctx, x.Db, user.ID, updateVals); errx != nil {
			slog.ErrorContext(ctx, "error updating user auth",
				logger.Err(errx))
			x.InternalServerError(w, "error updating user auth")
			return
		}
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) UpdateRemindersSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		x.BadFormDataError(w, err)
		return
	}

	changes := false
	enableReminders := r.PostFormValue("enable_reminders")
	notifThreshold := r.PostFormValue("notification_threshold")

	switch enableReminders {
	case "off":
		if user.Settings.EnableReminders {
			changes = true
			user.Settings.EnableReminders = false
		}
	case "on":
		if !user.Settings.EnableReminders {
			if !user.Verified {
				x.SessMgr.FlashAppend(ctx, "error", "Account must be verified before enabling reminder emails")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			user.Settings.EnableReminders = true
		}
	case "":
		// nothing
	default:
		x.SessMgr.FlashAppend(ctx, "error", "Bad value for reminders toggle")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	if notifThreshold != "" {
		v, err := strconv.ParseUint(notifThreshold, 10, 8)
		if err != nil {
			x.SessMgr.FlashAppend(ctx, "error", "Bad value for notification threshold")
			http.Redirect(w, r, "/settings", http.StatusSeeOther)
			return
		}
		val, err := model.ValidateReminderThresholdHours(v)
		if err != nil {
			x.SessMgr.FlashAppend(ctx, "error", "Bad value for notification threshold")
			http.Redirect(w, r, "/settings", http.StatusSeeOther)
			return
		}
		if user.Settings.ReminderThresholdHours != val {
			changes = true
			user.Settings.ReminderThresholdHours = val
		}
	}

	if !changes {
		x.SessMgr.FlashAppend(ctx, "error", "no changes made")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	if errx := service.UpdateUserSettings(
		ctx, x.Db, user.ID, &user.Settings); errx != nil {
		slog.ErrorContext(ctx, "error updating user settings",
			logger.Err(errx))
		x.InternalServerError(w, "error updating user settings")
		return
	}

	x.SessMgr.FlashAppend(ctx, "success", "update successful")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if errx := service.DeleteUser(ctx, x.Db, user.ID); errx != nil {
		x.DBError(w, errx)
		return
	}
	// destroy session
	err = x.SessMgr.Destroy(ctx)
	if err != nil {
		// not a fatal error (do not return 500 to user), since the user deleted
		// sucessfully already. just log the oddity
		slog.ErrorContext(ctx, "error destroying session",
			logger.Err(err))
	}
	if htmx.Hx(r).Request() {
		x.SessMgr.FlashAppend(ctx, "success", "Account deleted. Sorry to see you go.")
		w.Header().Add("HX-Location", "/login")
	}
	w.WriteHeader(200)
}
