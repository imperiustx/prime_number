package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/imperiustx/prime_number/cmd/server-api/internal/handlers"
	"github.com/imperiustx/prime_number/internal/platform/database"
	"github.com/imperiustx/prime_number/internal/schema"
)

// TestUsers runs a series of tests to exercise Product behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
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

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := UserTests{app: handlers.API(db, log)}

	t.Run("List", tests.List)
	t.Run("CreateRequiresFields", tests.CreateRequiresFields)
	t.Run("UserCRUD", tests.UserCRUD)
}

// UserTests holds methods for each user subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type UserTests struct {
	app http.Handler
}

func (u *UserTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/users", nil)
	resp := httptest.NewRecorder()

	u.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"id":           "5cf37266-3473-4006-984f-9325122678b7",
			"name":         "Admin Gopher",
			"email":        "admin@example.com",
			"roles":        []interface{}{"ADMIN", "USER"},
			"date_created": "2019-03-24T00:00:00Z",
			"date_updated": "2019-03-24T00:00:00Z",
		},
		{
			"id":           "45b5fbd3-755f-4379-8f07-a58d4a30fa2f",
			"name":         "User Gopher",
			"email":        "user@example.com",
			"roles":        []interface{}{"USER"},
			"date_created": "2019-03-24T00:00:00Z",
			"date_updated": "2019-03-24T00:00:00Z",
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (u *UserTests) CreateRequiresFields(t *testing.T) {
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/v1/users", body)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	u.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusBadRequest, resp.Code)
	}
}

func (u *UserTests) UserCRUD(t *testing.T) {
	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"user0","email":"user0@example.com","roles":["USER"],"password":"123456","password_confirm":"123456"}`)

		req := httptest.NewRequest("POST", "/v1/users", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		u.app.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("expected non-empty user id")
		}
		if created["date_created"] == "" || created["date_created"] == nil {
			t.Fatal("expected non-empty user date_created")
		}
		if created["date_updated"] == "" || created["date_updated"] == nil {
			t.Fatal("expected non-empty user date_updated")
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
			"name":         "user0",
			"email":        "user0@example.com",
			"roles":        created["roles"],
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/users/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		u.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched user should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved user should match created. Diff:\n%s", diff)
		}
	}

	{ // UPDATE
		body := strings.NewReader(`{"name":"new name"}`)
		url := fmt.Sprintf("/v1/users/%s", created["id"])
		req := httptest.NewRequest("PUT", url, body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		u.app.ServeHTTP(resp, req)

		if http.StatusNoContent != resp.Code {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp = httptest.NewRecorder()

		u.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var updated map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": updated["date_updated"],
			"name":         "new name",
			"email":        "user0@example.com",
			"roles":        created["roles"],
		}

		// Updated user should match the one we created.
		if diff := cmp.Diff(want, updated); diff != "" {
			t.Fatalf("Retrieved user should match created. Diff:\n%s", diff)
		}
	}

	{ // DELETE
		url := fmt.Sprintf("/v1/users/%s", created["id"])
		req := httptest.NewRequest("DELETE", url, nil)
		resp := httptest.NewRecorder()

		u.app.ServeHTTP(resp, req)

		if http.StatusNoContent != resp.Code {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp = httptest.NewRecorder()

		u.app.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusNotFound, resp.Code)
		}
	}
}
