package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
)

type mwContextKey string

func UserFromContext(ctx context.Context) (*model.User, error) {
	v, ok := ctx.Value(mwContextKey("user")).(*model.User)
	if !ok {
		return nil, fmt.Errorf("bad user context")
	}
	return v, nil
}

func IsLoggedIn(ctx context.Context) bool {
	v, ok := ctx.Value(mwContextKey("auth")).(bool)
	return ok && v
}

func LoadAuth(db *model.DB, sessMgr *session.SessionMgr) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if sessMgr.Exists(r.Context(), "user-id") {
				userId := sessMgr.Get(r.Context(), "user-id").(uint)
				user, err := model.GetUserById(db, r.Context(), userId)
				if err != nil {
					http.Error(w, "authorization failure", http.StatusInternalServerError)
					return
				}

				ctx = context.WithValue(ctx, mwContextKey("auth"), true)
				ctx = context.WithValue(ctx, mwContextKey("user"), user)
			} else {
				ctx = context.WithValue(ctx, mwContextKey("auth"), false)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := ctx.Value(mwContextKey("auth"))
		if v == nil || !v.(bool) {
			http.Error(w, "unauthorized, please login", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
