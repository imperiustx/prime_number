package handlers

import (
	"log"
	"net/http"

	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger) http.Handler {

	app := web.NewApp(log)

	u := Users{DB: db, Log: log}

	app.Handle(http.MethodGet, "/v1/users", u.List)
	app.Handle(http.MethodGet, "/v1/users/{id}", u.Retrieve)
	app.Handle(http.MethodPost, "/v1/users", u.Create)

	return app
}
