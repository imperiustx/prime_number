package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/imperiustx/prime_number/internal/platform/web"
	"github.com/imperiustx/prime_number/internal/primenumberrequest"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// PrimeNumberRequests defines all of the handlers related to prime number request. It holds the
// application state needed by the handler methods.
type PrimeNumberRequests struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all requests from the service layer and encodes them for the
// client response.
func (p *PrimeNumberRequests) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := primenumberrequest.List(ctx, p.DB)
	if err != nil {
		return errors.Wrap(err, "getting requests list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve finds a single user identified by an ID in the request URL.
func (p *PrimeNumberRequests) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	req, err := primenumberrequest.Retrieve(ctx, p.DB, id)
	if err != nil {
		switch err {
		case primenumberrequest.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case primenumberrequest.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting request %q", id)
		}
	}

	return web.Respond(ctx, w, req, http.StatusOK)
}

// Create decodes the body of a request to create a new request. The full
// request with generated fields is sent back in the response.
func (p *PrimeNumberRequests) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var nr primenumberrequest.NewRequest

	if err := web.Decode(r, &nr); err != nil {
		return errors.Wrap(err, "decoding new request")
	}

	req, err := primenumberrequest.Create(ctx, p.DB, claims, nr, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new requst")
	}

	return web.Respond(ctx, w, &req, http.StatusCreated)
}

// ListRequests gets all requests by an user from the service layer and encodes them for the
// client response.
func (p *PrimeNumberRequests) ListRequests(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "uid")
	list, err := primenumberrequest.ListRequests(ctx, p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting requests list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}
