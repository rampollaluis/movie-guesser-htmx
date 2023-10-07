package db

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Movie struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
	Genre string `json:"genre"`
}

type Movies []Movie

var allMovies map[string]string

func InitDb() error {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		return err
	}

    if err = createEmptyTablesIfNotExists(); err != nil {
		return err
	}

    loadMoviesFromFile()

	return err
}

func GetAllMovies() map[string]string {
    return allMovies
}

func createEmptyTablesIfNotExists() error {
	// Create a movies table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS movies (movie_id INTEGER PRIMARY KEY, name TEXT, emoji TEXT)")
	if err != nil {
		return err
	}
	statement.Exec()

	// Create a guesses table
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS guesses (movie_id INTEGER PRIMARY KEY, correct_count INTEGER, total_count INTEGER)")
	if err != nil {
		return err
	}
	statement.Exec()

	// Create a ratings table
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS ratings (movie_id INTEGER PRIMARY KEY, likes_count INTEGER, dislikes_count INTEGER)")
	if err != nil {
		return err
	}
	statement.Exec()

    statement.Close()

    return err
}

func loadMoviesFromFile() {
	jsonFile, err := os.Open("movies.json")
	if err != nil {
		log.Fatalf("Failed to open movies.json: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read from movies.json: %v", err)
	}

	movies := Movies{}

	if err := json.Unmarshal(byteValue, &movies); err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

    allMovies = make(map[string]string)
	for _, v := range movies {
		allMovies[v.Emoji] = v.Name
	}

	for _, movie := range movies {
		if _, err := FindMovieIdByEmoji(movie.Emoji); err != nil && err != sql.ErrNoRows {
			log.Printf("Failed to check for movie with emoji %s: %v", movie.Emoji, err)
			continue
		} else if err == nil {
			continue
		}

		err := InsertMovie(movie.Name, movie.Emoji)
		if err != nil {
			log.Printf("Failed to insert movie %s with emoji %s: %v", movie.Name, movie.Emoji, err)
		}
	}
    
}

func Close() {
	if db != nil {
		db.Close()
	}
}
