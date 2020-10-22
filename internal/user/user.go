package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// List gets all Users from the database.
func List(db *sqlx.DB) ([]User, error) {
	users := []User{}

	const q = `SELECT * FROM users`

	if err := db.Select(&users, q); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}

	return users, nil
}
