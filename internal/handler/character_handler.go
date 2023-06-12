package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"ricknmorty/internal/usecase/character"
	"strconv"
)

type CharacterHandler struct {
	FetchCharactersUseCase *character.FetchCharactersUsecase
}

func NewCharacterHandler(fetchCharactersUseCase *character.FetchCharactersUsecase) *CharacterHandler {
	return &CharacterHandler{FetchCharactersUseCase: fetchCharactersUseCase}
}

func (h *CharacterHandler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request from %s to %s\n", r.Method, r.RemoteAddr, r.RequestURI)

	queryParams := r.URL.Query()
	pageParam := queryParams.Get("page")
	page, err := strconv.Atoi(pageParam)

	if err != nil || page < 1 {
		page = 1
	}

	filters := make(map[string]string)
	for key, values := range queryParams {
		if len(values) > 0 && key != "page" {
			filters[key] = values[0]
		}
	}

	characters, err := h.FetchCharactersUseCase.Execute(filters, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
