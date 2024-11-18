package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Rebne/movie_generator/models"

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

	err = addMoviesFromJSON("./data/id_data.json")
	if err != nil {
		log.Printf("Warning: Could not load initial movie data: %v", err)
	}
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

func addMoviesFromJSON(filepath string) error {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var idData models.IDdata
	if err := json.Unmarshal(fileContent, &idData); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO movies (movie_id, title) 
		VALUES ($1, $2)
		ON CONFLICT (movie_id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, movie := range idData.IDs {
		_, err = stmt.Exec(movie.MovieID, movie.Title)
		if err != nil {
			return fmt.Errorf("error inserting movie %s: %v", movie.MovieID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
