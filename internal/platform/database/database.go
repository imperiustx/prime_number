// Package database helps with SQL database interactions.
package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

// Open knows how to open a database connection.
func Open() (*sqlx.DB, error) {
	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("root", "secret"),
		Host:     "localhost",
		Path:     "prime_number",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
