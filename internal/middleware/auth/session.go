package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/logger"
)

type mwContextKey string

func UserFromContext(ctx context.Context) (*model.User, error) {
	v, ok := ContextGet[*model.User](ctx, "user")
	if !ok || v == nil {
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
	v, ok := ContextGet[*model.User](ctx, "user")
	return ok && v != nil
}

type UserGetter interface {
	GetUserByID(context.Context, int) (*model.User, errs.Error)
}

type GetInter interface {
	GetInt(ctx context.Context, key string) int
}

func Load(userGetter UserGetter, intGetter GetInter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID := intGetter.GetInt(ctx, "user-id")
			if userID != 0 {
				user, err := userGetter.GetUserByID(ctx, userID)
				if err != nil {
					slog.InfoContext(ctx, "authorization failure", logger.Err(err))
					http.Error(w, "authorization failure", http.StatusUnauthorized)
					return
				}
				ctx = ContextSet(ctx, "user", user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !IsLoggedIn(ctx) {
			// if auth is required, and this is a get request,
			// redirect to login page and set "next=" query param
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
