// Package handlers provides the translation between the HTTP layer and application business logic.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
