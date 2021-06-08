package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/katier2d2/go-service/foundation/web"
)

func Panic(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the go stack tracw for the panic'd goroutine
					log.Printf("%s : PANIC   :\n%s", v.TraceID, debug.Stack())
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
