package handler

import (
	"net/http"

	"github.com/dropwhile/icbt/internal/app/middleware"
)

func (h *Handler) ShowIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try to get user from session
	user, err := middleware.UserFromContext(ctx)
	if err == nil && user != nil {
		// got a session, go to dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
