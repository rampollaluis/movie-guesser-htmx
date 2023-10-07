package db

func AddLike(movie_id int64) error {
	statement, err := db.Prepare("UPDATE ratings SET likes_count = likes_count + 1 WHERE movie_id = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}

func RemoveLike(movie_id int64) error {
	statement, err := db.Prepare("UPDATE ratings SET likes_count = likes_count - 1 WHERE movie_id = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}

func AddDislike(movie_id int64) error {
	statement, err := db.Prepare("UPDATE ratings SET dislikes_count = dislikes_count + 1 WHERE movie_id = ?")

	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}

func RemoveDislike(movie_id int64) error {
	statement, err := db.Prepare("UPDATE ratings SET dislikes_count = dislikes_count - 1 WHERE movie_id = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}

func SwapLikeForDislike(movie_id int64) error {
	statement, err := db.Prepare("UPDATE ratings SET likes_count = likes_count - 1, dislikes_count = dislikes_count + 1 WHERE movie_id = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}

func SwapDislikeForLike(movie_id int64) error {
	statement, err := db.Prepare("UPDATE ratings SET likes_count = likes_count + 1, dislikes_count = dislikes_count - 1 WHERE movie_id = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(movie_id)
    statement.Close()
	return err
}


