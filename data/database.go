package data

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

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

func initAdminAccount() {
	user := &Users{
		Email:    "admin@tnh.com",
		Name:     "Admin",
		Password: "1",
		Priority: math.MaxInt16,
	}

	user.Insert(user)
}

func initTable(db *sql.DB) error {
	// Create a table based on the Users struct
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT UNIQUE NOT NULL,
        name TEXT NOT NULL,
        password TEXT,
		authenticatorsecretkey TEXT,
		is2faenabled BOOLEAN,
		priority smallint NOT NULL,
        active BOOLEAN,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Table setup successfully")
	return nil
}
