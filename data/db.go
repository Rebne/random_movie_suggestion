package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Rebne/movie_generator/models"
	"github.com/Rebne/movie_generator/services"

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

func GetAllMoviesDB() (models.IDdata, error) {
	var result models.IDdata
	rows, err := db.Query("SELECT movie_id, title, index FROM movies")
	if err != nil {
		return models.IDdata{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id models.ID
		if err := rows.Scan(&id.MovieID, &id.Title, &id.Index); err != nil {
			return models.IDdata{}, err
		}
		result.IDs = append(result.IDs, id)
	}
	if err = rows.Err(); err != nil {
		return models.IDdata{}, err
	}
	result.Length, err = GetTableLengthDB()
	if err != nil {
		return models.IDdata{}, err
	}
	return result, nil
}

func GetTableLengthDB() (int, error) {
	var rowCount int
	err := db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&rowCount)
	if err != nil {
		return 0, err
	}
	return rowCount, nil
}

func GetAllMovieIdsDB() ([]string, error) {
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

func DeleteMovieDB(id string) error {
	result, err := db.Exec("DELETE FROM movies WHERE movie_id = $1 RETURNING *", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("ID %s does not exist", id)
	}
	return nil
}
func AddNewMovieDB(id string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE movie_id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if movie exists: %v", err)
	}
	if exists {
		return fmt.Errorf("movie with ID %s already exists", id)
	}
	movieData, err := services.FetchMovieData(id)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO movies (movie_id, title)
		VALUES ($1, $2)
	`, id, movieData.Title)
	if err != nil {
		return fmt.Errorf("error inserting movie: %v", err)
	}
	return nil
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
