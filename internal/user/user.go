package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

// Retrieve finds the user identified by a given ID.
func Retrieve(db *sqlx.DB, id string) (*User, error) {
	var u User

	const q = `SELECT * FROM users WHERE user_id = $1`

	if err := db.Get(&u, q, id); err != nil {
		return nil, errors.Wrap(err, "selecting single user")
	}

	return &u, nil
}

// Create adds a User to the database. It returns the created User with
// fields like ID and DateCreated populated..
func Create(db *sqlx.DB, nu NewUser, now time.Time) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `
		INSERT INTO users
		(user_id, name, email, roles, password_hash, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if _, err = db.Exec(q,
		u.ID, u.Name,
		u.Email, u.Roles, u.PasswordHash,
		u.DateCreated, u.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}

	return &u, nil
}
