package main

import (
	"main/store"
)

func main() {
	store := store.NewStore()
	err := store.Init()
	if err != nil {
		panic(err)
	}
	err = store.DB.Ping()
	if err != nil {
		panic(err)
	}
}
