package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

type FilterPair struct {
	filterValue string
	filterName  string
}

var db *sql.DB

func main() {
	var err error
	connStr := "user=vsevolodtrofimov dbname=characters sslmode=disable password="
	db, err = sql.Open("postgres", connStr)
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
	nameFilter := FilterPair{query.Get("name"), "name"}
	statusFilter := FilterPair{query.Get("status"), "status"}
	speciesFilter := FilterPair{query.Get("species"), "species"}
	typeFilter := FilterPair{query.Get("type"), "type"}
	genderFilter := FilterPair{query.Get("gender"), "gender"}

	filtersList := []FilterPair{nameFilter, statusFilter, speciesFilter, typeFilter, genderFilter}

	log.Printf("Name: %s, Status: %s, Species: %s, Type: %s, Gender: %s", nameFilter, statusFilter, speciesFilter, typeFilter, genderFilter)

	pageParam := query.Get("page")
	page, err := strconv.Atoi(pageParam)

	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 20
	offset := (page - 1) * pageSize

	whereFiltersClauses := []string{}
	params := []interface{}{}
	placeholderIndex := 1

	for _, item := range filtersList {
		if item.filterValue != "" {
			filterClause(item, &whereFiltersClauses, &params, &placeholderIndex)
		}
	}

	log.Printf("WhereFiltersClauses: %s", whereFiltersClauses)

	whereClause := ""
	if len(whereFiltersClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereFiltersClauses, " AND ")
	}

	log.Printf("Where clause: %s", whereClause)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM characters %s", whereClause)
	log.Printf("Count query: %s", countQuery)
	log.Printf("Params: %v", params)
	row := db.QueryRow(countQuery, params...)
	var totalCount int
	err = row.Scan(&totalCount)
	if err != nil {
		log.Printf("Error: %s", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	log.Printf("Total count: %d", totalCount)

	charactersQuery := fmt.Sprintf("SELECT * FROM characters %s LIMIT $%d OFFSET $%d", whereClause, placeholderIndex, placeholderIndex+1)
	log.Printf("Char query: %s", charactersQuery)
	params = append(params, pageSize, offset)
	rows, err := db.Query(charactersQuery, params...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var characters []Character
	for rows.Next() {
		var character Character
		var episode pq.StringArray

		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Status,
			&character.Species,
			&character.Type,
			&character.Gender,
			&character.Image,
			&character.Url,
			&character.Created,
			&character.Location.Name,
			&character.Location.Url,
			&character.Origin.Name,
			&character.Origin.Url,
			&episode,
		)

		character.Episode = episode

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		characters = append(characters, character)
	}

	// var totalCount int
	// err = db.QueryRow("SELECT count(*) FROM characters").Scan(&totalCount)
	// if err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	info := CharacterInfo{
		Count: totalCount,
		Pages: (totalCount + pageSize - 1) / pageSize,
	}

	if page > 1 {
		info.Prev = fmt.Sprintf("http://localhost:8080/api/character?page=%d", page-1)
	}

	if offset+pageSize < totalCount {
		info.Next = fmt.Sprintf("http://localhost:8080/api/character?page=%d", page+1)
	}

	response := APIResponse{
		Info:    info,
		Results: characters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func filterClause(filterPair FilterPair, whereQuery *[]string, params *[]interface{}, placeholderIndex *int) {
	log.Printf("Parms: %s, %s, WHERE: %s", filterPair.filterName, filterPair.filterValue, *whereQuery)
	*whereQuery = append(*whereQuery, fmt.Sprintf("%s ILIKE $%d", filterPair.filterName, *placeholderIndex))
	if filterPair.filterName != "gender" && filterPair.filterName != "species" {
		filterPair.filterValue = "%" + filterPair.filterValue + "%"
	}
	*params = append(*params, filterPair.filterValue)
	*placeholderIndex++
}
