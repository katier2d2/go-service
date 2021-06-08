package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/katier2d2/go-service/business/auth"
	"github.com/katier2d2/go-service/business/mid"
	"github.com/katier2d2/go-service/foundation/web"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, auth *auth.Auth) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panic(log))

	check := check{
		log: log,
	}

	app.Handle(http.MethodGet, "/readiness", check.readiness, mid.Authenticate(auth), mid.Authorize(log, "ADMIN"))

	return app
}
