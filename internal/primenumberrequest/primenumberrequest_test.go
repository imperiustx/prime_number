package primenumberrequest_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/imperiustx/prime_number/internal/platform/database"
	"github.com/imperiustx/prime_number/internal/primenumberrequest"
)

var result int64

func TestHighestPrimeNumber(t *testing.T) {
	r1 := primenumberrequest.HighestPrimeNumber(1234567)
	e1 := int64(1234547)

	if r1 != e1 {
		t.Errorf("expected '%d' but got '%d'", e1, r1)
	}

	r2 := primenumberrequest.HighestPrimeNumber(100)
	e2 := int64(97)

	if r2 != e2 {
		t.Errorf("expected '%d' but got '%d'", e2, r2)
	}

	r3 := primenumberrequest.HighestPrimeNumber(123)
	e3 := int64(113)

	if r3 != e3 {
		t.Errorf("expected '%d' but got '%d'", e3, r3)
	}

}

func BenchmarkHighestPrimeNumber1(b *testing.B) {
	var r int64
	for n := 0; n < b.N; n++ {
		r = primenumberrequest.HighestPrimeNumber(123)
	}

	result = r
}

func TestUsers(t *testing.T) {

	db, err := database.Open(database.Config{
		User:       "root",
		Password:   "secret",
		Host:       "localhost",
		Name:       "prime_number",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("connecting to db: %s", err)
	}
	defer db.Close()

	newR := primenumberrequest.NewRequest{
		SendNumber: 100,
	}

	now := time.Now()
	ctx := context.Background()

	claims := auth.NewClaims(
		"5cf37266-3473-4006-984f-9325122678b7", // This is just some random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

	r0, err := primenumberrequest.Create(ctx, db, claims, newR, now)
	if err != nil {
		t.Fatalf("creating request r0: %s", err)
	}

	r1, err := primenumberrequest.Retrieve(ctx, db, r0.ID)
	if err != nil {
		t.Fatalf("getting request r0: %s", err)
	}

	if diff := cmp.Diff(r1, r0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	saved, err := primenumberrequest.Retrieve(ctx, db, r0.ID)
	if err != nil {
		t.Fatalf("getting request r0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original request
	// and change just the fields we expect then diff it with what was saved.
	want := *r0
	want.ReceiveNumber = 97

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}

}
