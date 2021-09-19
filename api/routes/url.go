package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortner/domain"

	"github.com/rs/zerolog"
)

type urlHandlerService interface {
	CreateUrl(context.Context, string, string) (*domain.Url, error)
	CheckUrlExists(context.Context, string, string) (bool, error)
	GetUrl(context.Context, string) (*domain.Url, error)
}

type UrlHandler struct {
	Service urlHandlerService
	Log     *zerolog.Logger
}

type UserInput struct {
	Name       string `json:"name"`
	RedirectTo string `json:"redirectTo"`
}

func (uh *UrlHandler) CreateUrl(rw http.ResponseWriter, r *http.Request) {
	ui := &UserInput{}
	switch r.Header.Get("content-type") {
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(ui)
		if err != nil {
			uh.Log.Warn().Err(err).Msg("could not parse body")
			rw.WriteHeader(http.StatusNotAcceptable)
			rw.Write([]byte("could not parse body"))
			return
		}
	default:
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	ctx := r.Context()
	ok, err := uh.Service.CheckUrlExists(ctx, ui.Name, ui.RedirectTo)
	if err != nil {
		uh.Log.Fatal().Err(err).Msg("could not check existance of the url")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("url name already taken"))
		return
	}
	u, err := uh.Service.CreateUrl(ctx, ui.Name, ui.RedirectTo)
	if err != nil {
		uh.Log.Warn().Err(err).Msg("could not create url")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	json.NewEncoder(rw).Encode(u)
}
