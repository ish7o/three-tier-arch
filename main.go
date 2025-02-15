package main

import (
	"log"
	"net/http"
	"three-tier-arch/api"
	"three-tier-arch/store"
)


func main() {
	store := store.NewUserStore()
    handler := api.NewUserHandler(store)

	http.HandleFunc("/users", handler.HandleUsers)
	http.HandleFunc("/users/", handler.HandleUser)

	log.Fatal(http.ListenAndServe(":55005", nil))
}
