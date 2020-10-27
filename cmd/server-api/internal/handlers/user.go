// Package handlers provides the translation between the HTTP layer and application business logic.
package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/imperiustx/prime_number/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Users defines all of the handlers related to users. It holds the
// application state needed by the handler methods.
type Users struct {
	DB            *sqlx.DB
	Log           *log.Logger
	Authenticator *auth.Authenticator
}

// List gets all users from the service layer and encodes them for the
// client response.
func (u *Users) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := user.List(ctx, u.DB)
	if err != nil {
		return errors.Wrap(err, "getting user list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve finds a single user identified by an ID in the request URL.
func (u *Users) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	usr, err := user.Retrieve(ctx, u.DB, id)
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

	return web.Respond(ctx, w, usr, http.StatusOK)
}

// Create decodes the body of a request to create a new user. The full
// user with generated fields is sent back in the response.
func (u *Users) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nu user.NewUser

	if err := web.Decode(r, &nu); err != nil {
		return errors.Wrap(err, "decoding new user")
	}

	usr, err := user.Create(ctx, u.DB, nu, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new user")
	}

	return web.Respond(ctx, w, &usr, http.StatusCreated)
}

// Update decodes the body of a request to update an existing user. The ID
// of the user is part of the request URL.
func (u *Users) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update user.UpdateUser
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	if err := user.Update(ctx, u.DB, claims, id, update, time.Now()); err != nil {
		switch err {
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating user %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a single user identified by an ID in the request URL.
func (u *Users) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := user.Delete(ctx, u.DB, id); err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting user %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Token generates an authentication token for a user. The client must include
// an email and password for the request using HTTP Basic Auth. The user will
// be identified by email and authenticated by their password.
func (u *Users) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := user.Authenticate(ctx, u.DB, v.Start, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = u.Authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
