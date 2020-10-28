package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/imperiustx/prime_number/internal/mid"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(shutdown chan os.Signal, db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, log, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	{
		// Register health check handler. This route is not authenticated.
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		// Register user handlers.
		u := Users{DB: db, Log: log, Authenticator: authenticator}

		app.Handle(http.MethodGet, "/v1/users", u.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/users/{id}", u.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/users", u.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/v1/users/{id}", u.Update, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/v1/users/{id}", u.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

		// The token route can't be authenticated because they need this route to
		// get the token in the first place.
		app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	}

	{	
		// Register Prime Number Request handlers. Ensure all routes are authenticated.
		p := PrimeNumberRequests{DB: db, Log: log}

		app.Handle(http.MethodGet, "/v1/requests", p.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/requests/{id}", p.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/requests", p.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/requests/user/{uid}", p.ListRequests, mid.Authenticate(authenticator))
	}

	return app
}
