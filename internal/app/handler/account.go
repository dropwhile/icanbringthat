// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/samber/mo"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/htmx"
	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

func (x *Handler) AccountShowCreate(w http.ResponseWriter, r *http.Request) {
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
		"title":   "Create Account",
		"flashes": x.sessMgr.FlashPopAll(ctx),
		"next":    r.FormValue("next"),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "create-account-form.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) SettingsShow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	credentials, errx := x.svc.GetUserCredentialsByUser(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	apikey, errx := x.svc.GetApiKeyByUser(ctx, user.ID)
	if errx != nil {
		if errx.Code() != errs.NotFound {
			x.DBError(w, errx)
			return
		}
	}

	notifCount, errx := x.svc.GetNotificationsCount(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	// parse user-id url param
	tplVars := MapSA{
		"user":        user,
		"credentials": credentials,
		"apikey":      apikey,
		"title":       "Settings",
		"notifCount":  notifCount,
		"flashes":     x.sessMgr.FlashPopAll(ctx),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "show-settings.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) AccountCreate(w http.ResponseWriter, r *http.Request) {
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

	user, errx := x.svc.NewUser(ctx, email, name, []byte(passwd))
	if errx != nil {
		slog.ErrorContext(ctx, "error adding user", logger.Err(errx))
		x.BadRequestError(w, "error adding user")
		return
	}

	_, errx = x.svc.NewNotification(ctx, user.ID,
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
	err := x.sessMgr.RenewToken(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error renewing session token",
			logger.Err(err))
		x.InternalServerError(w, "Session Error")
		return
	}
	// Then make the privilege-level change.
	x.sessMgr.Put(r.Context(), "user-id", user.ID)
	x.sessMgr.FlashAppend(ctx, "success", "Account created. You are now logged in.")

	target := "/dashboard"
	if r.PostFormValue("next") != "" {
		target = r.FormValue("next")
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (x *Handler) SettingsUpdate(w http.ResponseWriter, r *http.Request) {
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
	updateVals := &service.UserUpdateValues{}

	email := r.PostFormValue("email")
	if email != "" && email != user.Email {
		changes = true
		updateVals.Email = mo.Some(email)
		updateVals.Verified = mo.Some(false)
		successMsgs = append(successMsgs, "Email update successfull")
	} else if email == user.Email {
		warnings = append(warnings, "Same Email specified was already present")
	}

	name := r.PostFormValue("name")
	if name != "" && name != user.Name {
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
			updateVals.PwUpdate = mo.Some(
				&service.PasswdUpdate{
					NewPass: []byte(newPasswd),
					OldPass: []byte(oldPasswd),
				},
			)
			changes = true
		}
	}

	if changes {
		errx := x.svc.UpdateUser(ctx, user, updateVals)
		if errx != nil {
			switch errx.Code() {
			case errs.InvalidArgument:
				arg := errx.Meta("argument")
				switch arg {
				case "OldPass":
					warnings = append(warnings, "Old Password invalid")
				case "Passwd":
					warnings = append(warnings, "'Password' was a bad value")
				default:
					warnings = append(warnings, errx.Msg())
				}
			default:
				slog.ErrorContext(ctx, "error updating user",
					logger.Err(errx))
				x.InternalServerError(w, "error updating user")
				return
			}
		}
	}
	if updateVals.PwUpdate.IsPresent() {
		successMsgs = append(successMsgs, "Password update successfull")
	}
	if len(warnings) > 0 {
		x.sessMgr.FlashAppend(ctx, "error", warnings...)
	} else if len(successMsgs) > 0 {
		x.sessMgr.FlashAppend(ctx, "success", successMsgs...)
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) SettingsAuthUpdate(w http.ResponseWriter, r *http.Request) {
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
	pkCount, errx := x.svc.GetUserCredentialCountByUser(ctx, user.ID)
	if errx != nil {
		x.DBError(w, errx)
		return
	}
	hasPasskeys := false
	if pkCount > 0 {
		hasPasskeys = true
	}

	updateVals := &service.UserUpdateValues{}
	switch authPW {
	case "off":
		if user.PWAuth {
			if !user.WebAuthn {
				x.sessMgr.FlashAppend(ctx, "error", "Account must be verified before enabling api access")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			updateVals.PWAuth = mo.Some(false)
		}
	case "on":
		if !user.PWAuth {
			changes = true
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
				x.sessMgr.FlashAppend(ctx, "error", "Refusing to disable passkey auth without alternative auth enabled")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			updateVals.WebAuthn = mo.Some(false)
		}
	case "on":
		if !user.WebAuthn {
			if !hasPasskeys {
				x.sessMgr.FlashAppend(ctx, "error",
					"Must have at least one passkey registered before enabling passkey auth")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			updateVals.WebAuthn = mo.Some(true)
		}
	case "":
		// nothing
	default:
		x.BadFormDataError(w, nil, "auth_passkeys")
		return
	}

	if !changes {
		x.sessMgr.FlashAppend(ctx, "error", "no changes made")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	errx = x.svc.UpdateUser(ctx, user, updateVals)
	if errx != nil {
		slog.ErrorContext(ctx, "error updating user auth",
			logger.Err(errx))
		x.InternalServerError(w, "error updating user auth")
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) SettingsAuthApiUpdate(w http.ResponseWriter, r *http.Request) {
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

	apiAccess := r.PostFormValue("api_access")
	rotateApiKey := r.PostFormValue("rotate_apikey")

	if apiAccess == "" && rotateApiKey == "" {
		x.BadFormDataError(w, nil, "bad params")
		return
	}

	changes := false
	updateVals := &service.UserUpdateValues{}
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
				x.sessMgr.FlashAppend(ctx, "error", "Refusing to enable api access without a verified account")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			updateVals.ApiAccess = mo.Some(true)
			_, errx := x.svc.NewApiKeyIfNotExists(ctx, user.ID)
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
		x.sessMgr.FlashAppend(ctx, "error", "no changes made")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	if rotateApiKey == "true" {
		if _, errx := x.svc.NewApiKey(ctx, user.ID); errx != nil {
			slog.ErrorContext(ctx, "error rotating api key",
				logger.Err(errx))
			x.InternalServerError(w, "error rotating api key")
			return
		}
	}

	if apiAccess != "" {
		if errx := x.svc.UpdateUser(ctx, user, updateVals); errx != nil {
			slog.ErrorContext(ctx, "error updating user auth",
				logger.Err(errx))
			x.InternalServerError(w, "error updating user auth")
			return
		}
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) SettingsRemindersUpdate(w http.ResponseWriter, r *http.Request) {
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
				x.sessMgr.FlashAppend(ctx, "error", "Account must be verified before enabling reminder emails")
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
				return
			}
			changes = true
			user.Settings.EnableReminders = true
		}
	case "":
		// nothing
	default:
		x.sessMgr.FlashAppend(ctx, "error", "Bad value for reminders toggle")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	if notifThreshold != "" {
		v, err := strconv.ParseUint(notifThreshold, 10, 8)
		if err != nil {
			x.sessMgr.FlashAppend(ctx, "error", "Bad value for notification threshold")
			http.Redirect(w, r, "/settings", http.StatusSeeOther)
			return
		}
		val, err := model.ValidateReminderThresholdHours(v)
		if err != nil {
			x.sessMgr.FlashAppend(ctx, "error", "Bad value for notification threshold")
			http.Redirect(w, r, "/settings", http.StatusSeeOther)
			return
		}
		if user.Settings.ReminderThresholdHours != val {
			changes = true
			user.Settings.ReminderThresholdHours = val
		}
	}

	if !changes {
		x.sessMgr.FlashAppend(ctx, "error", "no changes made")
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	if errx := x.svc.UpdateUserSettings(
		ctx, user.ID, &user.Settings); errx != nil {
		slog.ErrorContext(ctx, "error updating user settings",
			logger.Err(errx))
		x.InternalServerError(w, "error updating user settings")
		return
	}

	x.sessMgr.FlashAppend(ctx, "success", "update successful")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) AccountDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	if errx := x.svc.DeleteUser(ctx, user.ID); errx != nil {
		x.DBError(w, errx)
		return
	}
	// destroy session
	err = x.sessMgr.Destroy(ctx)
	if err != nil {
		// not a fatal error (do not return 500 to user), since the user deleted
		// sucessfully already. just log the oddity
		slog.ErrorContext(ctx, "error destroying session",
			logger.Err(err))
	}
	if htmx.Request(r).IsRequest() {
		x.sessMgr.FlashAppend(ctx, "success", "Account deleted. Sorry to see you go.")
		htmx.Response(w).HxLocation("/login")
	}
	w.WriteHeader(200)
}
