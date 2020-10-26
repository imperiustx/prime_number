// Package handlers provides the translation between the HTTP layer and application business logic.
package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/imperiustx/prime_number/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Users defines all of the handlers related to users. It holds the
// application state needed by the handler methods.
type Users struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all users from the service layer and encodes them for the
// client response.
func (u *Users) List(w http.ResponseWriter, r *http.Request) error {
	list, err := user.List(r.Context(), u.DB)
	if err != nil {
		return errors.Wrap(err, "getting user list")
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve finds a single user identified by an ID in the request URL.
func (u *Users) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	usr, err := user.Retrieve(r.Context(), u.DB, id)
	if err != nil {
		switch err {
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting user %q", id)
		}
	}

	return web.Respond(w, usr, http.StatusOK)
}

// Create decodes the body of a request to create a new user. The full
// user with generated fields is sent back in the response.
func (u *Users) Create(w http.ResponseWriter, r *http.Request) error {
	var nu user.NewUser

	if err := web.Decode(r, &nu); err != nil {
		return errors.Wrap(err, "decoding new user")
	}

	usr, err := user.Create(r.Context(), u.DB, nu, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new user")
	}

	return web.Respond(w, &usr, http.StatusCreated)
}
