package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"ricknmorty/internal/domain/service"
	"ricknmorty/internal/usecase/favourite"
	"strconv"
	"strings"
)

type FavouriteHandler struct {
	addFavoriteCharacter    *favourite.AddFavouriteCharacter
	removeFavoriteCharacter *favourite.RemoveFavouriteCharacter
	fetchFavoriteCharacters *favourite.FetchFavouriteCharacters
	tokenParser             *service.JwtTokenParser
}

func NewFavouriteHandler(
	addFavoriteCharacter *favourite.AddFavouriteCharacter,
	removeFavoriteCharacter *favourite.RemoveFavouriteCharacter,
	fetchFavoriteCharacters *favourite.FetchFavouriteCharacters,
	tokenParser *service.JwtTokenParser,
) *FavouriteHandler {
	return &FavouriteHandler{
		addFavoriteCharacter:    addFavoriteCharacter,
		removeFavoriteCharacter: removeFavoriteCharacter,
		fetchFavoriteCharacters: fetchFavoriteCharacters,
		tokenParser:             tokenParser,
	}
}

func (h *FavouriteHandler) getUserIdFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, http.ErrBodyNotAllowed
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.tokenParser.ParseToken(tokenStr)

	if err != nil {
		return 0, http.ErrBodyNotAllowed
	}

	return strconv.Atoi(claims.Subject)
}

func (h *FavouriteHandler) AddFavoriteCharacter(w http.ResponseWriter, r *http.Request) {
	log.Println("Got new character to favorites")

	userId, err := h.getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Invalid authorization", http.StatusUnauthorized)
		return
	}

	var data map[string]int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	characterId := data["characterId"]

	err = h.addFavoriteCharacter.Execute(userId, characterId)
	if err != nil {
		http.Error(w, "Unable to add favorite character", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FavouriteHandler) RemoveFavoriteCharacter(w http.ResponseWriter, r *http.Request) {
	log.Println("Got new character to remove from favorites")

	userId, err := h.getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Invalid authorization", http.StatusUnauthorized)
		return
	}

	var data map[string]int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	characterId := data["characterId"]

	err = h.removeFavoriteCharacter.Execute(userId, characterId)
	if err != nil {
		http.Error(w, "Unable to remove favorite character", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FavouriteHandler) FetchFavoriteCharacters(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching favorite characters")

	userId, err := h.getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Invalid authorization", http.StatusUnauthorized)
		return
	}

	characters, err := h.fetchFavoriteCharacters.Execute(userId)
	if err != nil {
		http.Error(w, "Unable to fetch favorite characters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
