package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"three-tier-arch/models"
	"three-tier-arch/store"
	// "three-tier-arch/models"
)


func isValidGroup(group string) bool {
	return group == "user" || group == "premium" || group == "admin"
}


func main() {
	userStore := store.NewUserStore()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			users := userStore.GetAllUsers()
            w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)

		case "POST":
			var input models.UserInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "invalid input sowwy", http.StatusBadRequest)
				return
			}
			user, err := userStore.CreateUser(input)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

            w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(user)

		default:
			http.Error(w, "method not awwowed oopsie", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
		if err != nil {
			http.Error(w, "bad id :c", http.StatusBadRequest)
            return
		}
		switch r.Method {
		case "GET":
            user, err := userStore.GetUser(id)
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

            user, err := userStore.UpdateUser(id, input)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(user)
		case "DELETE":
            err := userStore.DeleteUser(id)
            if err != nil {
                http.Error(w, err.Error(), http.StatusNotFound)
                return
            }

            w.WriteHeader(http.StatusNoContent)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":55005", nil))
}
