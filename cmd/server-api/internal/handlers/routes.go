package handlers

import (
	"log"
	"net/http"

	"github.com/imperiustx/prime_number/internal/mid"
	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger) http.Handler {

	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		u := Users{DB: db, Log: log}

		app.Handle(http.MethodGet, "/v1/users", u.List)
		app.Handle(http.MethodGet, "/v1/users/{id}", u.Retrieve)
		app.Handle(http.MethodPost, "/v1/users", u.Create)
		app.Handle(http.MethodPut, "/v1/users/{id}", u.Update)
		app.Handle(http.MethodDelete, "/v1/users/{id}", u.Delete)
	}

	{
		p := PrimeNumberRequests{DB: db, Log: log}

		app.Handle(http.MethodGet, "/v1/requests", p.List)
		app.Handle(http.MethodGet, "/v1/requests/{id}", p.Retrieve)
		app.Handle(http.MethodPost, "/v1/requests", p.Create)
		app.Handle(http.MethodGet, "/v1/requests/user/{uid}", p.ListRequests)
	}

	return app
}
