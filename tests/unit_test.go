package tests

import (
	"testing"
	"three-tier-arch/models"
	"three-tier-arch/store"
	"time"
)

// Rawr
var (
	testInputA = models.UserInput{FirstName: "Ada", LastName: "Lovelace", BirthYear: 1815, Group: "premium"}
	testInputB = models.UserInput{FirstName: "Brian", LastName: "Kernighan", BirthYear: 1942, Group: "user"}
	testInputC = models.UserInput{FirstName: "Clifford", LastName: "Cocks", BirthYear: 1950, Group: "admin"}
)

func TestUserStore(t *testing.T) {

	// GET /users & POST /users
	t.Run("GetAllUsers", func(t *testing.T) {
        store := store.NewUserStore()

		users := store.GetAllUsers()
		if len(users) != 0 {
			t.Errorf("Expected empty store, got %d users", len(users))
		}

		store.CreateUser(testInputA)
		store.CreateUser(testInputB)

		users = store.GetAllUsers()
		if len(users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(users))
		}

		store.DeleteUser(1)
		users = store.GetAllUsers()
		if len(users) != 1 {
			t.Errorf("Expected 1 user, got %d", len(users))
		}

		store.CreateUser(testInputC)

		users = store.GetAllUsers()
		if len(users) != 2 {
			t.Errorf("Expected 2 user, got %d", len(users))
		}
	})

	// GET /users/<id>
	t.Run("GetUser", func(t *testing.T) {
		store := store.NewUserStore()
		_, err := store.GetUser(1)
		if err == nil {
			t.Error("Expected error for non-existent user")
		}

		created, _ := store.CreateUser(testInputA)
		user, err := store.GetUser(created.Id)
		if err != nil {
			t.Errorf("Failed to get created user: %v", err)
		}
		if user.FirstName != testInputA.FirstName {
			t.Errorf("Expected first name %s; got %s", testInputA.FirstName, user.FirstName)
		}
		if user.LastName != testInputA.LastName {
			t.Errorf("Expected last name %s; got %s", testInputA.LastName, user.LastName)
		}
		if user.Group != testInputA.Group {
			t.Errorf("Expected group %s; got %s", testInputA.Group, user.Group)
		}
		if user.BirthYear != testInputA.BirthYear {
			t.Errorf("Expected birth year %d; got %d", testInputA.BirthYear, user.BirthYear)
		}
		if user.Age != time.Now().Year()-testInputA.BirthYear {
			t.Errorf("Expected age %d; got %d", time.Now().Year()-testInputA.BirthYear, user.Age)
		}
	})

	// PATCH /users/<id>
	t.Run("UpdateUser", func(t *testing.T) {
		store := store.NewUserStore()

		user, _ := store.CreateUser(testInputB)
		store.UpdateUser(user.Id, testInputC)
		userToo, err := store.GetUser(user.Id)
		if err != nil {
			t.Fatalf("Expected user got error: %v", err)
		}
		if testInputC.FirstName != userToo.FirstName {
			t.Fatalf("User data not updated")
		}
	})

	// DELETE /users/<id>
	t.Run("DeleteUser", func(t *testing.T) {
		store := store.NewUserStore()

		user, _ := store.CreateUser(testInputC)
		err := store.DeleteUser(user.Id)
		if err != nil {
			t.Errorf("Failed to delete existing user: %v", err)
		}

		_, err = store.GetUser(user.Id)
		if err == nil {
			t.Error("Expected error when getting deleted user")
		}

		err = store.DeleteUser(2137)
		if err == nil {
			t.Error("Expected error when deleting non-existent user")
		}
	})
}
