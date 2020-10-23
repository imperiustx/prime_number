package user_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/imperiustx/prime_number/internal/platform/database/databasetest"
	"github.com/imperiustx/prime_number/internal/schema"
	"github.com/imperiustx/prime_number/internal/user"
)

func TestUsers(t *testing.T) {

	_, db, teardown := databasetest.NewUnit(t)

	defer teardown()

	newU := user.NewUser{
		Name:            "Thor",
		Email:           "thor@marvel.com",
		Roles:           []string{"USER"},
		Password:        "123456",
		PasswordConfirm: "123456",
	}

	now := time.Now()

	u0, err := user.Create(db, newU, now)
	if err != nil {
		t.Fatalf("creating user u0: %s", err)
	}

	u1, err := user.Retrieve(db, u0.ID)
	if err != nil {
		t.Fatalf("getting user u0: %s", err)
	}

	if diff := cmp.Diff(u1, u0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

}

func TestUserList(t *testing.T) {
	_, db, teardown := databasetest.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := user.List(db)
	if err != nil {
		t.Fatalf("listing users: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected user list size %v, got %v", exp, got)
	}
}
