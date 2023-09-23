package zhandler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
)

func (z *ZHandler) ShowIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := auth.UserFromContext(ctx)
	if err == nil && user != nil {
		// got a session, go to dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
