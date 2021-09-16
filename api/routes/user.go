package routes

import (
	"main/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/user", controllers.RegisterUser).Methods(http.MethodPost)
}
