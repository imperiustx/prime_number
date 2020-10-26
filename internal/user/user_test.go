package user_test

import (
	"context"
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
	ctx := context.Background()

	u0, err := user.Create(ctx, db, newU, now)
	if err != nil {
		t.Fatalf("creating user u0: %s", err)
	}

	u1, err := user.Retrieve(ctx, db, u0.ID)
	if err != nil {
		t.Fatalf("getting user u0: %s", err)
	}

	if diff := cmp.Diff(u1, u0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}

	update := user.UpdateUser{
		Name: databasetest.StringPointer("Thor Odinson"),
	}
	updatedTime := time.Now()

	if err := user.Update(ctx, db, u0.ID, update, updatedTime); err != nil {
		t.Fatalf("updating user u0: %s", err)
	}

	saved, err := user.Retrieve(ctx, db, u0.ID)
	if err != nil {
		t.Fatalf("getting user u0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original user
	// and change just the fields we expect then diff it with what was saved.
	want := *u0
	want.Name = "Thor Odinson"
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}

	if err := user.Delete(ctx, db, u0.ID); err != nil {
		t.Fatalf("deleting product: %v", err)
	}

	_, err = user.Retrieve(ctx, db, u0.ID)
	if err == nil {
		t.Fatalf("should not be able to retrieve deleted user")
	}

}

func TestUserList(t *testing.T) {
	_, db, teardown := databasetest.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := user.List(context.Background(), db)
	if err != nil {
		t.Fatalf("listing users: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected user list size %v, got %v", exp, got)
	}
}
