package main

import (
	"emoji-movies-htmx/db"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
	}
	return err
}

type Data struct {
	Movies    []string
	Selection string
}

var allMovies map[string]string

func Hello(c echo.Context) error {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	values := make([]string, 0)
	keys := make([]string, 0)
	for k, v := range allMovies {
		values = append(values, v)
		keys = append(keys, k)
	}
	data := Data{
		Movies:    values,
		Selection: keys[r.Intn(len(keys))],
	}

	return c.Render(http.StatusOK, "hello", data)
}

type OutcomeData struct {
    Emoji string
    IsCorrect     bool
    Correct_count int64
    Total_count   int64
}

func (data OutcomeData) PercentageRight() string {
    return fmt.Sprintf("%.2f", float32(data.Correct_count) / float32(data.Total_count) * 100)
}

func Guess(c echo.Context) error {
	isCorrect := false
	emoji := c.FormValue("emojiValue")

	if strings.ToLower(allMovies[emoji]) == strings.ToLower(c.FormValue("answer")) {
		isCorrect = true
	}

	movie_id, err := db.FindMovieIdByEmoji(emoji)
	if err != nil {
		log.Fatal(err)
	}

	err = db.RecordGuess(movie_id, isCorrect)
	if err != nil {
		log.Fatal(err)
	}


	correct_count, total_count, err := db.GetCorrectAndTotalGuessesByMovieId(movie_id)
	if err != nil {
		log.Fatal(err)
	}

	outcomeData := OutcomeData {
        emoji,
		isCorrect,
		correct_count,
		total_count,
	}

	return c.Render(http.StatusOK, "outcome", outcomeData)
}

func VoteUp(c echo.Context) error {
	log.Println("thumbs up clicked")

	emoji := c.FormValue("emojiValue")
    movie_id, err := db.FindMovieIdByEmoji(emoji)
    if err != nil {
        log.Fatal(err)
    }

	if thumbsUpAlreadySelected, _ := strconv.ParseBool(c.FormValue("thumbsUpSelected")); thumbsUpAlreadySelected {
		log.Println("reversing thumbs up")
        err = db.RemoveLike(movie_id)
	} else if thumbsDownSelected, _ := strconv.ParseBool(c.FormValue("thumbsDownSelected")); thumbsDownSelected {
		log.Println("changing thumbs down for thumbs up")
        err = db.SwapDislikeForLike(movie_id)
	} else {
		log.Println("adding thumbs up to question")
        err = db.AddLike(movie_id)
	}

    if err != nil {
        log.Fatal(err)
    }

	return err
}

func VoteDown(c echo.Context) error {
	log.Println("thumbs down clicked")

	emoji := c.FormValue("emojiValue")
    movie_id, err := db.FindMovieIdByEmoji(emoji)
    if err != nil {
        log.Fatal(err)
    }

	if thumbsDownAlreadySelected, _ := strconv.ParseBool(c.FormValue("thumbsDownSelected")); thumbsDownAlreadySelected {
		log.Println("reversing thumbs down")
        err = db.RemoveDislike(movie_id)
	} else if thumbsUpSelected, _ := strconv.ParseBool(c.FormValue("thumbsUpSelected")); thumbsUpSelected {
		log.Println("changing thumbs up for thumbs down")
        err = db.SwapLikeForDislike(movie_id)
	} else {
		log.Println("adding thumbs down to question")
        err = db.AddDislike(movie_id)
	}

    if err != nil {
        log.Fatal(err)
    }

	return nil
}

func main() {
	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}
    defer db.Close()

    allMovies = db.GetAllMovies()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Static("/static", "static")
	e.GET("/", Hello)
	e.POST("/guess", Guess)
	e.POST("/vote-down", VoteDown)
	e.POST("/vote-up", VoteUp)
	e.Logger.Fatal(e.Start(":3000"))
}

