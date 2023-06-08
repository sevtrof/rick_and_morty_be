package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	var characters []Character

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

		characters = append(characters, apiResponse.Results...)

		if apiResponse.Info.Next == "" {
			break
		}
		apiURL = apiResponse.Info.Next
	}

	data, err := json.Marshal(characters)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("characters.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
