package mid

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/katier2d2/go-service/business/auth"
	"github.com/katier2d2/go-service/foundation/web"
)

var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)

func Authenticate(a *auth.Auth) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			// validate
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = context.WithValue(ctx, auth.Key, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

func Authorize(log *log.Logger, roles ...string) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from ctx")
			}

			if !claims.Authorize(roles...) {
				log.Printf("mid: authorize: claims: %v expecting: %v", claims.Roles, roles)
				return ErrForbidden
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
