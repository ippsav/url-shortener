package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"url-shortner/domain"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type urlHandlerService interface {
	CreateUrl(context.Context, string, string) (*domain.Url, error)
	CheckUrlExists(context.Context, string, string) (bool, error)
	GetUrlByID(context.Context, int64) (*domain.Url, error)
	GetUrls(context.Context, int, time.Time) ([]domain.Url, error)
	GetUrlByName(context.Context, string) (*domain.Url, error)
}

type UrlHandler struct {
	Service urlHandlerService
	Log     *zerolog.Logger
}

type urlJsonFormat struct {
	Name       string `json:"name"`
	RedirectTo string `json:"redirectTo"`
}
type urlQueryArgs struct {
	CreatedAt time.Time `json:"createdAt"`
	Limit     int       `json:"limit"`
}

func (uh *UrlHandler) CreateUrl(rw http.ResponseWriter, r *http.Request) {
	ui := &urlJsonFormat{}
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

func (uh *UrlHandler) GetUrl(rw http.ResponseWriter, r *http.Request) {
	uName := chi.URLParam(r, "name")
	_ = uName
}

func (uh *UrlHandler) GetUrls(rw http.ResponseWriter, r *http.Request) {
	query := &urlQueryArgs{}
	switch r.Header.Get("content-type") {
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(query)
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
	urls, err := uh.Service.GetUrls(ctx, query.Limit, query.CreatedAt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(urls)
}
