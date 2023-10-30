package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type movie struct {
	Title   string
	Year    int
	Runtime int
	Genres  []string
}

func (i movie) escapedTitle() string {
	return strings.ReplaceAll(i.Title, "'", "''")
}

func (i movie) joinGenres() string {
	joined := strings.Join(i.Genres, "\",\"")
	return fmt.Sprintf("'{\"%s\"}'", joined)
}

//go:embed movies.json
var f embed.FS

func main() {
	file, err := f.Open("movies.json")
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(file)

	if _, err = dec.Token(); err != nil {
		log.Fatal(err)
	}

	values := []string{}

	for dec.More() {
		var m movie
		if err = dec.Decode(&m); err != nil {
			log.Fatal(err)
		}

		v := fmt.Sprintf("('%s', %d, %d, %v)", m.escapedTitle(), m.Year, m.Runtime, m.joinGenres())
		values = append(values, v)
	}

	fmt.Printf(
		"insert into movies(title, year, runtime, genres) values %s;",
		strings.Join(values, ",\n"),
	)

	if _, err = dec.Token(); err != nil {
		log.Fatal(err)
	}
}
