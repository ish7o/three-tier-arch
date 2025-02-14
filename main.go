package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

type UserStore struct {
    sync.Mutex
    users map[int]User
    nextId int
}

type UserInput struct {
    FirstName string `json:"firstName"`
    LastName string `json:"lastname"`
    BirthYear int `json:"birthYear"`
    Group string `json:"group"`
}

type User struct {
    Id int `json:"id"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastname"`
    BirthYear int `json:"birthYear"`
    Age int `json:"age"`
    Group string `json:"group"`
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

func (us *UserStore) AddUser(input UserInput) (*User, error) {
    if input.FirstName == "" || input.LastName == "" || !isValidGroup(input.Group) {
        return nil, errors.New("Input data not cute enough :333")
    }

    age := time.Now().Year() - input.BirthYear
    us.Lock()
    defer us.Unlock()

    id := us.NewId()

    user := User{
        Id: id,
        FirstName: input.FirstName,
        LastName: input.LastName,
        BirthYear: age,
        Group: input.Group,
    }
    us.users[id] = user
    return &user, nil
}



func (us *UserStore) GetAllUsers() []User {
    us.Lock()
    defer us.Unlock()
    v := getMapValues(us.users)
    return v
}

func NewUserStore() *UserStore {
    return &UserStore{
        sync.Mutex{},
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
                json.NewEncoder(w).Encode(users)
            case "POST":
                var input UserInput
                if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
                    http.Error(w, "invalid input sowwy", http.StatusBadRequest)
                    return
                }
                user, err := userStore.AddUser(input)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusBadRequest)
                    return
                }

                w.WriteHeader(http.StatusCreated)
                json.NewEncoder(w).Encode(user)

            default:
                http.Error(w, "method not awwowed oopsie", http.StatusMethodNotAllowed)
        }
    })
    http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        // id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
        // if err != nil {
        //     http.Error(w, "bad id :c", http.StatusBadRequest)
        // }
        // switch r.Method {
        // case "GET":
        // case "PATCH":
        // case "DELETE":
        // }
    })

    log.Fatal(http.ListenAndServe(":55005", nil))
}
