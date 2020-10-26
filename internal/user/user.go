package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("user not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// List gets all Users from the database.
func List(ctx context.Context, db *sqlx.DB) ([]User, error) {
	users := []User{}

	const q = `SELECT * FROM users`

	if err := db.SelectContext(ctx, &users, q); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}

	return users, nil
}

// Retrieve finds the user identified by a given ID.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var u User

	const q = `SELECT * FROM users WHERE user_id = $1`

	if err := db.GetContext(ctx, &u, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting single user")
	}

	return &u, nil
}

// Create adds a User to the database. It returns the created User with
// fields like ID and DateCreated populated..
func Create(ctx context.Context, db *sqlx.DB, nu NewUser, now time.Time) (*User, error) {
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

	if _, err = db.ExecContext(ctx, q,
		u.ID, u.Name,
		u.Email, u.Roles, u.PasswordHash,
		u.DateCreated, u.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting user")
	}

	return &u, nil
}
