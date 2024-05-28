package main

import (
	"database/sql"
	"erp-user-service/data"
	"erp-user-service/data/utils"
	"erp-user-service/factory"
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
	DB           *sql.DB
	Models       *data.Models
	Utils        *utils.AppUtilities
	ErrorFactory *factory.ErrorFactory
}

const webPort = "80"

func main() {
	db := connectDB()

	if db == nil {
		log.Panic("Could not connect to postgres")
	}

	log.Println("Connected to postgres")

	app := Config{
		DB:           db,
		Models:       data.New(db),
		Utils:        utils.New(),
		ErrorFactory: &factory.ErrorFactory{},
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

		if err == nil {
			initTable(connection)
			return connection
		}

		count++

		if count > 10 {
			return nil
		}
		time.Sleep(2 * time.Second)
		continue
	}
}

func initTable(db *sql.DB) error {
	// Create a table based on the Users struct
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT UNIQUE NOT NULL,
        name TEXT,
        password TEXT NOT NULL,
        active BOOLEAN,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );`

	// Execute the SQL statement
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Table created successfully")
	return nil
}
