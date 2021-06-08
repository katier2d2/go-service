package mid

import (
	"context"
	"expvar"
	"net/http"
	"runtime"

	"github.com/katier2d2/go-service/foundation/web"
)

var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := handler(ctx, w, r)

			// increment the request
			m.req.Add(1)

			//Update the count for the number of active goroutines every 100 requests
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			//Increment the errors counter if an error occurred on this request
			if err != nil {
				m.err.Add(1)
			}

			return err
		}
		return h
	}
	return m
}
