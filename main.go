package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
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

type Movie struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
	Genre string `json:"genre"`
}

type Movies []Movie

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

func Guess(c echo.Context) error {
	// TODO: save to db and output correct guesses
	correct := false

	if strings.ToLower(allMovies[c.FormValue("emojiValue")]) == strings.ToLower(c.FormValue("answer")) {
		correct = true
	}

	return c.Render(http.StatusOK, "outcome", correct)
}

func VoteUp(c echo.Context) error {
	log.Println("thumbs up clicked")
	if thumbsUpAlreadySelected, _ := strconv.ParseBool(c.FormValue("thumbsUpSelected")); thumbsUpAlreadySelected {
		log.Println("reversing thumbs up")
	} else if thumbsDownSelected, _ := strconv.ParseBool(c.FormValue("thumbsDownSelected")); thumbsDownSelected {
		log.Println("changing thumbs down for thumbs up")
	} else {
		log.Println("adding thumbs up to question")
	}
	return nil
}

func VoteDown(c echo.Context) error {
	log.Println("thumbs down clicked")
	if thumbsDownAlreadySelected, _ := strconv.ParseBool(c.FormValue("thumbsDownSelected")); thumbsDownAlreadySelected {
		log.Println("reversing thumbs down")
	} else if thumbsUpSelected, _ := strconv.ParseBool(c.FormValue("thumbsUpSelected")); thumbsUpSelected {
		log.Println("changing thumbs up for thumbs down")
	} else {
		log.Println("adding thumbs down to question")
	}

	return nil
}

func main() {
	loadAllMovies()

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

func loadAllMovies() {
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
}
