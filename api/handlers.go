package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"three-tier-arch/models"
	"three-tier-arch/store"
)

type UserHandler struct {
	store *store.UserStore
}

func NewUserHandler(store *store.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (uh *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users := uh.store.GetAllUsers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case "POST":
		var input models.UserInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid input sowwy", http.StatusBadRequest)
			return
		}
		user, err := uh.store.CreateUser(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "method not awwowed", http.StatusMethodNotAllowed)
	}
}

func (uh *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
	if err != nil {
		http.Error(w, "bad id :c", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		user, err := uh.store.GetUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&user)

	case "PATCH":
		var input models.UserInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

		user, err := uh.store.UpdateUser(id, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	case "DELETE":
		err := uh.store.DeleteUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
