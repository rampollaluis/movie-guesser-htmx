package db

import (
	"database/sql"
	"log"
)

func RecordGuess(movie_id int64, isCorrect bool) error {
	var statement *sql.Stmt
	var err error

	if isCorrect {
		statement, err = db.Prepare("UPDATE guesses SET correct_count = correct_count + 1, total_count = total_count + 1 WHERE movie_id = ?")
		log.Printf("recording correct guess for %v\n", movie_id)
	} else {
		statement, err = db.Prepare("UPDATE guesses SET total_count = total_count + 1 WHERE movie_id = ?")
		log.Printf("recording incorrect guess for %v\n", movie_id)
	}

	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)

    statement.Close()
	return err
}

func GetCorrectAndTotalGuessesByMovieId(movie_id int64) (int64, int64, error) {
	log.Printf("getting guess stats for movie %v\n", movie_id)
	row := db.QueryRow("SELECT correct_count, total_count FROM guesses WHERE movie_id = ?", movie_id)

	var correct_count, total_count int64
	err := row.Scan(&correct_count, &total_count)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no match")
			return 0, 0, nil
		}
		return 0, 0, err
	}

	return correct_count, total_count, nil
}
