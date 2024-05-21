package main

import (
	"database/sql"
	"erp-user-service/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

const webPort = "80"

func main() {
	db := connectDB()

	if db == nil {
		log.Panic("Could not connect to postgres")
	}

	log.Println("Connected to postgres")

	app := Config{
		DB:     db,
		Models: data.New(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
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

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	count := 0

	for {
		connection, err := openDB(dsn)

		if err != nil {
			return connection
		}

		count++

		if count > 10 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
}
