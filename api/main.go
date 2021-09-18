package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"
	"url-shortner/routes"
	"url-shortner/service/user"
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

	// TODO: having the connection url as an environment variable insteading of setting it in the main func

	if err := os.Setenv("DB_URI", "root:password@(localhost:7200)/db?parseTime=true"); err != nil {
		log.Fatal().Err(err).Msg("Could not set the db uri")
	}
	//Connection to database
	db_uri := os.Getenv("DB_URI")
	if db_uri == "" {
		log.Fatal().Msg("Could not find the db_uri in the environment variables")
	}
	db, err := sql.Open("mysql", db_uri)
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
	log.Info().Msg("Users table is set")

	//Setting chi router
	r := chi.NewMux()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middlewareLogger(log))

	us := &user.Service{Store: store}
	uh := &routes.UserHandler{Service: us, Log: &log}

	//Api Status
	r.Get("/status", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Server running"))
	})
	//User Routes
	r.Post("/users", uh.CreateUser)
	r.Get("/users", uh.FindUser)
	// Serving mux router
	log.Info().Msg("server running on port 7000")
	if err := http.ListenAndServe(":7000", r); err != nil {
		log.Fatal().Err(err)
	}
}

func middlewareLogger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "logger", log)
			r = r.WithContext(ctx)

			log.Info().
				Str("remote_ip", r.RemoteAddr).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Request Received")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
