package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/suhel-kap/auth-service/cmd/data"
)

const (
	WEB_PORT = "80"
)

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the auth service")

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to Postgres")
	}

	// Setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not available, sleeping")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return conn
		}

		if counts > 10 {
			log.Println("Could not connect to Postgres")
			return nil
		}

		log.Println("Sleeping for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
