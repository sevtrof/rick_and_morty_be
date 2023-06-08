package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

var characters []Character

func main() {

	data, err := os.ReadFile("characters.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &characters)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/character", charactersHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request from %s to %s\n", r.Method, r.RemoteAddr, r.RequestURI)

	query := r.URL.Query()
	pageParam := query.Get("page")
	page, err := strconv.Atoi(pageParam)

	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 20
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(characters) {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	if end > len(characters) {
		end = len(characters)
	}

	info := CharacterInfo{
		Count: len(characters),
		Pages: (len(characters) + pageSize - 1) / pageSize,
	}

	if page > 1 {
		info.Prev = fmt.Sprintf("http://localhost:8080/api/character?page=%d", page-1)
	}

	if end < len(characters) {
		info.Next = fmt.Sprintf("http://localhost:8080/api/character?page=%d", page+1)
	}

	response := APIResponse{
		Info:    info,
		Results: characters[start:end],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
