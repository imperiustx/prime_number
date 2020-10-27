package primenumberrequest

import "time"

// PrimeNumberRequest is the information of which user request and receive which number
type PrimeNumberRequest struct {
	ID            string    `db:"request_id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	SendNumber    int64     `db:"send_number" json:"send_number"`
	ReceiveNumber int64     `db:"receive_number" json:"receive_number"`
	DateCreated   time.Time `db:"date_created" json:"date_created"`
}

// NewRequest is what we require from clients when they send a new request
type NewRequest struct {
	SendNumber int64  `json:"send_number" validate:"gte=3"`
}
