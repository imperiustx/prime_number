package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/imperiustx/prime_number/internal/platform/auth"
	"github.com/imperiustx/prime_number/internal/platform/database"
	"github.com/imperiustx/prime_number/internal/schema"
	"github.com/imperiustx/prime_number/internal/user"
	_ "github.com/lib/pq" // The database driver in use.
)

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

	newU := user.NewUser{
		Name:            "Thor",
		Email:           "thor@marvel.com",
		Roles:           []string{"USER"},
		Password:        "123456",
		PasswordConfirm: "123456",
	}

	now := time.Now()
	ctx := context.Background()

	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e", // This is just some random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

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
		Name: StringPointer("Thor Odinson"),
	}
	updatedTime := time.Now().UTC()

	if err := user.Update(ctx, db, claims, u0.ID, update, updatedTime); err != nil {
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

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}
