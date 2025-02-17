package main

import (
	"log"
	"net/http"
	"three-tier-arch/api"
	"three-tier-arch/service"
	"three-tier-arch/store"
)


func main() {
	store := store.NewUserStore() // DATA tier
    service := service.NewUserService(store) // Business~ logic (application tier)
    handler := api.NewUserHandler(service) // PRESENTATION TIER~

	http.HandleFunc("/users", handler.HandleUsers)
	http.HandleFunc("/users/", handler.HandleUser)

	log.Fatal(http.ListenAndServe(":55005", nil))
}
