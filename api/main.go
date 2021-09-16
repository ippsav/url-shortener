package main

import (
	"main/routes"
	"main/store"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Init Store
	store := store.NewStore()
	err := store.Init()
	if err != nil {
		panic(err)
	}
	// Checking server status
	err = store.DB.Ping()
	if err != nil {
		panic(err)
	}
	muxRouter := mux.NewRouter()
	routes.RegisterUserRoutes(muxRouter)

	if err := http.ListenAndServe(":7000", muxRouter); err != nil {
		panic(err)
	}
}
