package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"three-tier-arch/api"
	"three-tier-arch/models"
	"three-tier-arch/service"
	"three-tier-arch/store"
)

func TestIntegration(t *testing.T) {
	store := store.NewUserStore()
    service := service.NewUserService(store)
    handler := api.NewUserHandler(service)

    t.Run("Create and Get user and users", func(t *testing.T) {
        // make a user
        p, _ := json.Marshal(testInputA)
        postReq := httptest.NewRequest("POST", "/users", bytes.NewBuffer(p))
        postRec := httptest.NewRecorder()
        handler.HandleUsers(postRec, postReq)

        if postRec.Code != http.StatusCreated {
            t.Fatalf("Expected 201, got %d", postRec.Code)
        }

        var createdUser models.User
        json.NewDecoder(postRec.Body).Decode(&createdUser)


        // Get all users (only one ^)
        getReq := httptest.NewRequest("GET", "/users", nil)
        getRec := httptest.NewRecorder()
        handler.HandleUsers(getRec, getReq)

        if getRec.Code != http.StatusOK {
            t.Fatalf("Expected 200, got %d", getRec.Code)
        }

        var users []models.User
        err := json.NewDecoder(getRec.Body).Decode(&users)
        if err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if len(users) != 1 {
            t.Fatalf("Expected 1 user, got %d", len(users))
        }

        // get single user from id
        getReq = httptest.NewRequest("GET", "/users/" + string(strconv.Itoa(createdUser.Id)), nil)
        getRec = httptest.NewRecorder()

        handler.HandleUser(getRec, getReq)

        if getRec.Code != http.StatusOK {
            t.Fatalf("Expected 200, got %d", getRec.Code)
        }
    })

    t.Run("Update user", func(t *testing.T) {
        p, _ := json.Marshal(testInputB)
        postReq := httptest.NewRequest("POST", "/users", bytes.NewBuffer(p))
        postRec := httptest.NewRecorder()
        handler.HandleUsers(postRec, postReq)

        var user models.User
        json.NewDecoder(postRec.Body).Decode(&user)

        p, _ = json.Marshal(testInputC)
        patchReq := httptest.NewRequest("PATCH", "/users/"+strconv.Itoa(user.Id), bytes.NewBuffer(p))
        patchRec := httptest.NewRecorder()
        handler.HandleUser(patchRec, patchReq)

        if patchRec.Code != http.StatusOK {
            t.Fatalf("Expected 200, got %d", patchRec.Code)
        }
    })

    t.Run("Delete User", func(t *testing.T) {
        payload, _ := json.Marshal(testInputA)
        postReq := httptest.NewRequest("POST", "/users", bytes.NewBuffer(payload))
        postRec := httptest.NewRecorder()
        handler.HandleUsers(postRec, postReq)

        var user models.User
        json.NewDecoder(postRec.Body).Decode(&user)

        deleteReq := httptest.NewRequest("DELETE", "/users/"+strconv.Itoa(user.Id), nil)
        deleteRec := httptest.NewRecorder()
        handler.HandleUser(deleteRec, deleteReq)

        if deleteRec.Code != http.StatusNoContent {
            t.Fatalf("Expected 204, got %d", deleteRec.Code)
        }

        getReq := httptest.NewRequest("GET", "/users/"+strconv.Itoa(user.Id), nil)
        getRec := httptest.NewRecorder()
        handler.HandleUser(getRec, getReq)

        if getRec.Code != http.StatusNotFound {
            t.Fatalf("Expected 404, got %d", getRec.Code)
        }
    })
}
