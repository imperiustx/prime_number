// Package handlers provides the translation between the HTTP layer and application business logic.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/imperiustx/prime_number/internal/user"
	"github.com/jmoiron/sqlx"
)

// Users defines all of the handlers related to users. It holds the
// application state needed by the handler methods.
type Users struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all users from the service layer and encodes them for the
// client response.
func (u *Users) List(w http.ResponseWriter, r *http.Request) {
	list, err := user.List(u.DB)
	if err != nil {
		u.Log.Printf("error: listing users: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		u.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		u.Log.Println("error writing result", err)
	}
}

// Retrieve finds a single user identified by an ID in the request URL.
func (u *Users) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	usr, err := user.Retrieve(u.DB, id)
	if err != nil {
		u.Log.Println("getting Users", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(usr)
	if err != nil {
		u.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		u.Log.Println("error writing result", err)
	}
}

// Create decodes the body of a request to create a new user. The full
// user with generated fields is sent back in the response.
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var nu user.NewUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		u.Log.Println("decoding user", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prod, err := user.Create(u.DB, nu, time.Now())
	if err != nil {
		u.Log.Println("creating user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(prod)
	if err != nil {
		u.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		u.Log.Println("error writing result", err)
	}
}
