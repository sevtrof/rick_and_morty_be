package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"ricknmorty/internal/usecase/favourite"
	"strconv"

	"github.com/golang-jwt/jwt"
)

type FavouriteHandler struct {
	addFavoriteCharacter    *favourite.AddFavouriteCharacter
	removeFavoriteCharacter *favourite.RemoveFavouriteCharacter
	fetchFavoriteCharacters *favourite.FetchFavouriteCharacters
}

func NewFavouriteHandler(
	addFavoriteCharacter *favourite.AddFavouriteCharacter,
	removeFavoriteCharacter *favourite.RemoveFavouriteCharacter,
	fetchFavoriteCharacters *favourite.FetchFavouriteCharacters) *FavouriteHandler {
	return &FavouriteHandler{
		addFavoriteCharacter:    addFavoriteCharacter,
		removeFavoriteCharacter: removeFavoriteCharacter,
		fetchFavoriteCharacters: fetchFavoriteCharacters,
	}
}

func (h *FavouriteHandler) AddFavoriteCharacter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got new character to favorites")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := authHeader[7:]
	claims := &jwt.StandardClaims{}

	log.Printf("Parsing token")
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	log.Printf("Checking token: %v", token)
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("Token's fine")

	log.Printf("Getting userId")
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	log.Printf("Got userId: %d", userId)

	var data map[string]int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Decoding: data: %v, err: %v", data, err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	characterId := data["characterId"]

	log.Printf("UserId and characterId: %d, %d", userId, characterId)

	err = h.addFavoriteCharacter.Execute(userId, characterId)
	if err != nil {
		http.Error(w, "unable to add favorite character", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FavouriteHandler) RemoveFavoriteCharacter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got new character to remove from favorites")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := authHeader[7:]
	claims := &jwt.StandardClaims{}

	log.Printf("Parsing token")
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	log.Printf("Checking token: %v", token)
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("Token's fine")

	log.Printf("Getting userId")
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	log.Printf("Got userId: %d", userId)

	var data map[string]int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Decoding: data: %v, err: %v", data, err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	characterId := data["characterId"]

	log.Printf("UserId and characterId: %d, %d", userId, characterId)

	err = h.removeFavoriteCharacter.Execute(userId, characterId)
	if err != nil {
		http.Error(w, "unable to remove favorite character", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FavouriteHandler) FetchFavoriteCharacters(w http.ResponseWriter, r *http.Request) {
	log.Printf("Fetching favorite characters")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := authHeader[7:]
	claims := &jwt.StandardClaims{}

	log.Printf("Parsing token")
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	log.Printf("Checking token: %v", token)
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("Token's fine")

	log.Printf("Getting userId")
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	log.Printf("Got userId: %d", userId)

	characters, err := h.fetchFavoriteCharacters.Execute(userId)
	if err != nil {
		log.Printf("Fetching favourtie chars: %v, error: %v", characters, err)
		http.Error(w, "unable to fetch favorite characters", http.StatusInternalServerError)
		return
	}

	log.Printf("Got chars: %v", characters)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}
