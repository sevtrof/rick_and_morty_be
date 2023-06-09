package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type CharacterInfo struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type Character struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Species  string `json:"species"`
	Type     string `json:"type"`
	Gender   string `json:"gender"`
	Image    string `json:"image"`
	Url      string `json:"url"`
	Created  string `json:"created"`
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Origin struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"origin"`
	Episode []string `json:"episode"`
}

type APIResponse struct {
	Info    CharacterInfo `json:"info"`
	Results []Character   `json:"results"`
}

func main() {
	apiURL := "https://rickandmortyapi.com/api/character"

	connStr := "user= dbname= sslmode=disable password="
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS characters (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		status VARCHAR(255),
		species VARCHAR(255),
		type VARCHAR(255),
		gender VARCHAR(255),
		image VARCHAR(255),
		url VARCHAR(255),
		created TIMESTAMP,
		location_name VARCHAR(255),
		location_url VARCHAR(255),
		origin_name VARCHAR(255),
		origin_url VARCHAR(255),
		episode TEXT
	);`)

	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Fatal(err)
		}

		var apiResponse APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResponse)
		resp.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		for _, character := range apiResponse.Results {
			_, err := db.Exec(`INSERT INTO characters (name, status, species, type, gender, image, url, created, location_name, location_url, origin_name, origin_url, episode) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
				character.Name, character.Status, character.Species, character.Type, character.Gender, character.Image, character.Url, character.Created,
				character.Location.Name, character.Location.Url, character.Origin.Name, character.Origin.Url, pq.Array(character.Episode))
			if err != nil {
				log.Fatal(err)
			}
		}

		if apiResponse.Info.Next == "" {
			break
		}
		apiURL = apiResponse.Info.Next
	}

	fmt.Println("Done")
}
