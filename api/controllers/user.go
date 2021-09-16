package controllers

import (
	"encoding/json"
	"main/models"
	"main/service"
	"net/http"
)

var (
	userService = service.NewUserService()
)

type RegisterOrLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(rw http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	userInput := RegisterOrLoginInput{}
	json.NewDecoder(r.Body).Decode(&userInput)
	user.Email = userInput.Email
	user.PasswordHash = userInput.Password
	err := userService.Validate(user)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		rw.Write([]byte(err.Error()))
	}
	user, err = userService.CreateUser(user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(user)
}
