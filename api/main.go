package main

import (
	"main/store"
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

}
