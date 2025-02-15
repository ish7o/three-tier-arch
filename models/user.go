package models

import (
	"errors"
	"strings"
)

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

// I hate OOP and yet I come up with such nonsenses. (sorry my dear haskell)
func (u UserInput) IsValid() error {
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

func isValidGroup(group string) bool {
	return group == "user" || group == "premium" || group == "admin"
}

