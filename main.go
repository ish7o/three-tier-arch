package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type UserStore struct {
	sync.RWMutex
	users  map[int]User
	nextId int
}

func (us *UserStore) GetUser(id int) (*User, error) {
    us.RLock()
    defer us.RUnlock()
    user, exists := us.users[id]
    if !exists {
        return nil, fmt.Errorf("user with id %d not found (yet)", id)
    }
    return &user, nil
}

func (us *UserStore) UpdateUser(id int, input UserInput) (*User, error) {
    if err := input.isValid(); err != nil {
        return nil, fmt.Errorf("invalid user input: %w", err)
    }

    us.Lock()
    defer us.Unlock()

    if _, exists := us.users[id]; !exists {
        return nil, fmt.Errorf("user with id %d not found (yet)", id)
    }

	age := time.Now().Year() - input.BirthYear

    updatedUser := &User{
        Id: id,
        FirstName: input.FirstName,
        LastName: input.LastName,
        BirthYear: input.BirthYear,
        Group: input.Group,
        Age: age,
    }

    us.users[id] = *updatedUser
    return updatedUser, nil
}

func (us *UserStore) DeleteUser(id int) error {
    us.Lock()
    defer us.Unlock()

    if _, exists := us.users[id]; !exists {
        return fmt.Errorf("user with id %d not found (yet)", id)
    }

    delete(us.users, id)
    return nil
}

// I hate OOP and yet I come up with such nonsenses. (sorry my dear haskell)
func (u UserInput) isValid() error {
    if strings.TrimSpace(u.FirstName) == "" {
        return errors.New("First name is literally nothing!")
    }
    if strings.TrimSpace(u.LastName) == "" {
        return errors.New("Your last name... may just be nothing, sorry")
    }
    if !isValidGroup(u.Group) {
        return errors.New("your group should be either 'user', 'premium' or 'admin' *rolls eyes*")
    }
    return nil
}

func (us *UserStore) CreateUser(input UserInput) (*User, error) {
    if err := input.isValid(); err != nil {
        return nil, fmt.Errorf("invalid user input: %w", err)
    }

	age := time.Now().Year() - input.BirthYear
	us.Lock()
	defer us.Unlock()

    id := us.NewId()

	user := User{
		Id:        id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthYear: input.BirthYear,
		Group:     input.Group,
        Age:       age,
	}

	us.users[id] = user
	return &user, nil
}




type UserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastname"`
	BirthYear int    `json:"birthYear"`
	Group     string `json:"group"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastname"`
	BirthYear int    `json:"birthYear"`
	Age       int    `json:"age"`
	Group     string `json:"group"`
}

func (us *UserStore) NewId() int {
	us.nextId += 1
	return us.nextId
}

func isValidGroup(group string) bool {
	return group == "user" || group == "premium" || group == "admin"
}

// Brilliant function, literally co cool, and map key has to be of interface `Comparable` lmao
func getMapValues[T comparable, U any](m map[T]U) []U {
	s := make([]U, 0, len(m))
	for _, u := range m {
		s = append(s, u)
	}

	return s
}

func (us *UserStore) GetAllUsers() []User {
	us.RLock()
	defer us.RUnlock()
	v := getMapValues(us.users)

	return v
}

func NewUserStore() *UserStore {
	return &UserStore{
		sync.RWMutex{},
		make(map[int]User),
		0,
	}
}

func main() {
	userStore := NewUserStore()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			users := userStore.GetAllUsers()
            w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)

		case "POST":
			var input UserInput
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
            var input UserInput
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
