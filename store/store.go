package store

import (
	"fmt"
	"sync"
	"three-tier-arch/models"
	"time"
)

type UserStore struct {
	sync.RWMutex
	users  map[int]models.User
	nextId int
}

func NewUserStore() *UserStore {
	return &UserStore{
		sync.RWMutex{},
		make(map[int]models.User),
		0,
	}
}

func (us *UserStore) GetUser(id int) (*models.User, error) {
    us.RLock()
    defer us.RUnlock()
    user, exists := us.users[id]
    if !exists {
        return nil, fmt.Errorf("user with id %d not found (yet)", id)
    }
    return &user, nil
}

func (us *UserStore) UpdateUser(id int, input models.UserInput) (*models.User, error) {
    us.Lock()
    defer us.Unlock()

    if _, exists := us.users[id]; !exists {
        return nil, fmt.Errorf("user with id %d not found (yet)", id)
    }

	age := time.Now().Year() - input.BirthYear

    updatedUser := &models.User{
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

func (us *UserStore) CreateUser(input models.UserInput) (*models.User, error) {
	age := time.Now().Year() - input.BirthYear
	us.Lock()
	defer us.Unlock()

    id := us.NewId()

	user := models.User{
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

func (us *UserStore) NewId() int {
	us.nextId += 1
	return us.nextId
}

func (us *UserStore) GetAllUsers() []models.User {
	us.RLock()
	defer us.RUnlock()
	v := getMapValues(us.users)

	return v
}

// Brilliant function, literally co cool, and map key has to be of interface `Comparable` lmao
func getMapValues[T comparable, U any](m map[T]U) []U {
	s := make([]U, 0, len(m))
	for _, u := range m {
		s = append(s, u)
	}

	return s
}
