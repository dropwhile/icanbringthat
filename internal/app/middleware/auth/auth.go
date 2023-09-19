package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/rs/zerolog/log"
)

type mwContextKey string

func UserFromContext(ctx context.Context) (*model.User, error) {
	v, ok := ContextGet[*model.User](ctx, "user")
	if !ok {
		return nil, fmt.Errorf("bad context value for user")
	}
	return v, nil
}

func ContextGet[T any](ctx context.Context, key string) (T, bool) {
	v, ok := ctx.Value(mwContextKey(key)).(T)
	return v, ok
}

func ContextSet(ctx context.Context, key string, value any) context.Context {
	ctx = context.WithValue(ctx, mwContextKey(key), value)
	return ctx
}

func IsLoggedIn(ctx context.Context) bool {
	v, ok := ContextGet[bool](ctx, "auth")
	return ok && v
}

func Load(db model.PgxHandle, sessMgr *session.SessionMgr) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if sessMgr.Exists(r.Context(), "user-id") {
				userId := sessMgr.Get(r.Context(), "user-id").(int)
				user, err := model.GetUserById(r.Context(), db, userId)
				if err != nil {
					log.Err(err).Msg("authorization failure")
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

func Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := ctx.Value(mwContextKey("auth"))
		if v == nil || !v.(bool) {
			if r.Method == http.MethodGet ||
				r.Method == http.MethodHead {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "unauthorized, please log in", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
