package main

import (
	"context"
	"database/sql"
	"os"
	"time"

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

}
