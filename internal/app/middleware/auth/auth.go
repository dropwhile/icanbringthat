package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/modelx"
	"github.com/dropwhile/icbt/internal/session"
)

type mwContextKey string

func UserFromContext(ctx context.Context) (*modelx.User, error) {
	v, ok := ContextGet[*modelx.User](ctx, "user")
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

func Load(db modelx.PgxHandle, sessMgr *session.SessionMgr) func(next http.Handler) http.Handler {
	query := modelx.New(db)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if sessMgr.Exists(r.Context(), "user-id") {
				userID := sessMgr.Get(r.Context(), "user-id").(int32)
				user, err := query.GetUserById(r.Context(), userID)
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
				target := "/login"
				if !strings.HasPrefix(r.URL.Path, "/login") {
					q := url.Values{}
					q.Set("next", r.URL.Path)
					target = strings.Join([]string{target, q.Encode()}, "?")
				}
				http.Redirect(w, r, target, http.StatusSeeOther)
				return
			}
			http.Error(w, "unauthorized, please log in", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
