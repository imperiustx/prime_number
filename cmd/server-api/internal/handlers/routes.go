package handlers

import (
	"log"
	"net/http"

	"github.com/imperiustx/prime_number/internal/mid"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		u := Users{DB: db, Log: log, Authenticator: authenticator}

		app.Handle(http.MethodGet, "/v1/users", u.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/users/{id}", u.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/users", u.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/v1/users/{id}", u.Update, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/v1/users/{id}", u.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
		app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	}

	{
		p := PrimeNumberRequests{DB: db, Log: log}

		app.Handle(http.MethodGet, "/v1/requests", p.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/requests/{id}", p.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/requests", p.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/requests/user/{uid}", p.ListRequests, mid.Authenticate(authenticator))
	}

	return app
}
