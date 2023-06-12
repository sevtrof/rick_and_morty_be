package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ricknmorty/internal/handler"
	"ricknmorty/internal/repository"
	"ricknmorty/internal/usecase/character"
)

func main() {
	connStr := "user= dbname= sslmode=disable password="
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	characterRepo := repository.NewCharacterRepository(db)
	fetchCharacters := character.NewFetchCharacters(characterRepo)
	characterHandler := handler.NewCharacterHandler(fetchCharacters)

	http.HandleFunc("/api/character", characterHandler.GetCharacters)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
