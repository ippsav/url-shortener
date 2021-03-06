package main

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"net/http"
	"os"
	"time"
	"url-shortner/routes"
	"url-shortner/service/url"
	"url-shortner/service/user"
	credis "url-shortner/store/cache"
	"url-shortner/store/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

func main() {
	stdout := zerolog.NewConsoleWriter()
	log := zerolog.New(stdout)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// TODO: having the connection url as an environment variable instead of setting it in the main func

	if err := os.Setenv("DB_URI", "root:password@(localhost:7200)/db?parseTime=true"); err != nil {
		log.Fatal().Err(err).Msg("Could not set the db uri")
	}
	//Connection to database
	dbUri := os.Getenv("DB_URI")
	if dbUri == "" {
		log.Fatal().Msg("Could not find the dbUri in the environment variables")
	}
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not ping the database")
	}
	log.Info().Msg("Connected to database")

	//Creating mysql store
	store := &mysql.Store{DB: db}

	// Create users table
	if err := store.CreateUsersTable(ctx); err != nil {
		log.Fatal().Err(err).Msg("Could not create the users table")
	}
	log.Info().Msg("users table is set")

	// Create urls table
	if err := store.CreateUrlsTable(ctx); err != nil {
		log.Fatal().Err(err).Msg("Could not create the urls table")
	}
	log.Info().Msg("urls table is set")

	// create redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	client.FlushDB(ctx)
	redisCache := credis.NewRedisCache(client)
	//Setting chi router
	r := chi.NewMux()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middlewareLogger(log))
	//Services
	us := &user.Service{Store: store}
	urs := &url.Service{Store: store, Cache: redisCache}
	//Handlers
	uh := &routes.UserHandler{Service: us, Log: &log}
	urh := &routes.UrlHandler{Service: urs, Log: &log}

	//Api Status
	r.Get("/status", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Server running"))
	})
	//User Routes
	r.Post("/users", uh.CreateUser)
	r.Get("/users", uh.FindUser)
	r.Route("/urls", func(r chi.Router) {
		r.Use(authMiddleware(us))
		r.Post("/", urh.CreateUrl)
		r.Get("/", urh.GetUrls)
		r.Get("/{name}", urh.GetUrl)
	})
	// Serving mux router
	log.Info().Msg("server running on port 7000")
	if err := http.ListenAndServe(":7000", r); err != nil {
		log.Fatal().Err(err)
	}
}

func authMiddleware(us *user.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("qid")
			if err != nil {
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			ctx := r.Context()
			ok, err := us.Store.CheckUserExistsWithID(ctx, cookie.Value)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !ok {
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			ctx = context.WithValue(ctx, "userID", cookie.Value)
			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}

func middlewareLogger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "logger", log)
			r = r.WithContext(ctx)

			log.Info().
				Str("remote_ip", r.RemoteAddr).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Request Received")

			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}
