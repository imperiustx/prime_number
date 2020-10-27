package primenumberrequest

import (
	"context"
	"database/sql"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Request is requested but does not exist.
	ErrNotFound = errors.New("request not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// List gets all Requests from the database.
func List(ctx context.Context, db *sqlx.DB) ([]PrimeNumberRequest, error) {
	requests := []PrimeNumberRequest{}

	const q = `SELECT * FROM prime_number_requests`

	if err := db.SelectContext(ctx, &requests, q); err != nil {
		return nil, errors.Wrap(err, "selecting requests")
	}

	return requests, nil
}

// Retrieve finds the request identified by a given ID.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*PrimeNumberRequest, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p PrimeNumberRequest

	const q = `SELECT * FROM prime_number_requests WHERE request_id = $1`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting single request")
	}

	return &p, nil
}

// Create adds a Request to the database. It returns the created Request with
// fields like ID and DateCreated populated..
func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, nr NewRequest, now time.Time) (*PrimeNumberRequest, error) {

	p := PrimeNumberRequest{
		ID:            uuid.New().String(),
		UserID:        user.Subject,
		SendNumber:    nr.SendNumber,
		ReceiveNumber: highestPrimeNumber(nr.SendNumber),
		DateCreated:   now.UTC(),
	}

	const q = `
		INSERT INTO prime_number_requests
		(request_id, user_id, send_number, receive_number, date_created)
		VALUES ($1, $2, $3, $4, $5)`

	if _, err := db.ExecContext(ctx, q,
		p.ID, p.UserID,
		p.SendNumber, p.ReceiveNumber,
		p.DateCreated); err != nil {
		return nil, errors.Wrap(err, "inserting request")
	}

	return &p, nil
}

func highestPrimeNumber(num int64) int64 {
Prime:
	for {
		switch {
		case big.NewInt(num).ProbablyPrime(0):
			break Prime
		default:
			num--
		}
	}
	return num
}

// ListRequests gives all Requests of an User.
func ListRequests(ctx context.Context, db *sqlx.DB, userID string) ([]PrimeNumberRequest, error) {
	requests := []PrimeNumberRequest{}

	const q = `SELECT * FROM prime_number_requests WHERE user_id = $1`
	if err := db.SelectContext(ctx, &requests, q, userID); err != nil {
		return nil, errors.Wrap(err, "selecting requests by user")
	}

	return requests, nil
}
