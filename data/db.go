package data

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB
var connStr string

type MovieID struct {
	ID string `db:"id"`
}

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	connStr = os.Getenv("DATABSE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to Postgres database")
}

func InitDB() *sql.DB {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		movie_id TEXT NOT NULL UNIQUE
	)`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	return db
}
