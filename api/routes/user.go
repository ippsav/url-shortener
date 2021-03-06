package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortner/domain"

	"github.com/rs/zerolog"
)

type userHandlerService interface {
	CreateUser(context.Context, string, string) (*domain.User, error)
	FindUser(context.Context, string, string) (*domain.User, error)
	CheckUserExists(context.Context, string) (bool, error)
}

type UserHandler struct {
	Service userHandlerService
	Log     *zerolog.Logger
}

type registerLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInput := &registerLoginInput{}
	switch r.Header.Get("content-type") {
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(userInput)
		if err != nil {
			rw.WriteHeader(http.StatusNotAcceptable)
			rw.Write([]byte("error parsing body"))
			return
		}
	default:
		rw.WriteHeader(http.StatusNotAcceptable)
		rw.Write([]byte("Invalid content type"))
		return
	}
	ok, err := us.Service.CheckUserExists(ctx, userInput.Email)
	if err != nil {
		us.Log.Warn().Err(err).Msg("create user Error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok {
		us.Log.Info().Msg("user exists in database")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("email already taken"))
		return
	}

	u, err := us.Service.CreateUser(ctx, userInput.Email, userInput.Password)
	if err != nil {
		us.Log.Warn().Err(err).Msg("create user Error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(u)
}

func (us *UserHandler) FindUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInput := &registerLoginInput{}
	switch r.Header.Get("content-type") {
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(userInput)
		if err != nil {
			rw.WriteHeader(http.StatusNotAcceptable)
			rw.Write([]byte("error parsing body"))
			return
		}
	default:
		rw.WriteHeader(http.StatusNotAcceptable)
		rw.Write([]byte("Invalid content type"))
		return
	}
	u, err := us.Service.FindUser(ctx, userInput.Email, userInput.Password)
	if err != nil {
		us.Log.Warn().Err(err).Msg("could not find user")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:     "qid",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Value:    u.ID,
		MaxAge:   6000,
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(u)
}
