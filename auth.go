package jasmine

import (
	"context"
	"errors"
	"net/http"
)

type contextKey int

var authCtxKey contextKey = 0

type Context struct {
	UserName string
}

func NewContext(ctx context.Context, a *Context) context.Context {
	return context.WithValue(ctx, authCtxKey, a)
}

func FromContext(ctx context.Context) (*Context, bool) {
	a, ok := ctx.Value(authCtxKey).(*Context)
	return a, ok
}

type ProviderMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

var authError = errors.New("auth context error")

type Auth struct{}

func (u *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.Header().Add("WWW-Authenticate", `Basic realm="Please provide a password", charset="UTF-8"`)
			http.Error(w, "HTTP auth is required", http.StatusUnauthorized)
			return
		}
		username, accessToken, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Authorization invalid", http.StatusUnauthorized)
			return
		}

		// TODO: verify accessToken later
		_ = accessToken

		r = r.WithContext(NewContext(r.Context(), &Context{UserName: username}))
		r.BasicAuth()
		next.ServeHTTP(w, r)
	})
}

func (u *Auth) CurrentUserPrincipal(ctx context.Context) (string, error) {
	authCtx, ok := FromContext(ctx)
	if !ok {
		return "", authError
	}
	return "/" + authCtx.UserName + "/", nil
}
