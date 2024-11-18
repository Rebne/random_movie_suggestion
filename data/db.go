package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB
var connStr string

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	connStr = os.Getenv("DATABSE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging:", err)
	}
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS movies (
		index SERIAL PRIMARY KEY,
		movie_id TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL UNIQUE
	)`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	fmt.Println("Successfully connected to Postgres database")
}

func GetTableLengthDB() (int, error) {
	var rowCount int
	err := db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&rowCount)
	if err != nil {
		return 0, err
	}
	return rowCount, nil
}

func GetAllMoviesDB() ([]string, error) {
	var movieIDs []string
	rows, err := db.Query("SELECT movie_id FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		movieIDs = append(movieIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movieIDs, nil
}
