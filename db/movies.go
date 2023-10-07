package db

import (
	"database/sql"
	"log"
)


func InsertMovie(name string, emoji string) error {
	log.Printf("inserting movie %s", name)
	statement, err := db.Prepare("INSERT INTO movies (name, emoji) VALUES (?, ?)")
	if err != nil {
		return err
	}
	result, err := statement.Exec(name, emoji)
	if err != nil {
		return err
	}

	movie_id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	statement, err = db.Prepare("INSERT INTO guesses (movie_id, correct_count, total_count) VALUES (?, 0, 0)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(movie_id)

	statement, err = db.Prepare("INSERT INTO ratings (movie_id, likes_count, dislikes_count) VALUES (?, 0, 0)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}

func FindMovieIdByEmoji(emoji string) (int64, error) {
	log.Printf("searching for movie with emoji %v\n", emoji)
	row := db.QueryRow("SELECT movie_id FROM movies WHERE emoji = ?", emoji)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no match")
			return 0, err
		}
		return 0, err
	}

	return id, nil
}
