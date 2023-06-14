package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ricknmorty/internal/handler"
	"ricknmorty/internal/repository"
	"ricknmorty/internal/usecase/character"
	"ricknmorty/internal/usecase/favourite"
	"ricknmorty/internal/usecase/user"
)

func main() {
	connStr := "user= dbname= sslmode=disable password="
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Character
	characterRepo := repository.NewCharacterRepository(db)
	fetchCharacters := character.NewFetchCharacters(characterRepo)
	characterHandler := handler.NewCharacterHandler(fetchCharacters)

	// User
	userRepo := repository.NewUserRepository(db)
	loginUserUseCase := user.NewLoginUserUseCase(userRepo)
	logoutUserUseCase := user.NewLogoutUserUseCase()
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(loginUserUseCase, logoutUserUseCase, registerUserUseCase)

	// Favourite
	addFavoriteCharacter := favourite.NewAddFavouriteCharacter(userRepo)
	removeFavoriteCharacter := favourite.NewRemoveFavouriteCharacter(userRepo)
	fetchFavoriteCharacters := favourite.NewFetchFavouriteCharacters(userRepo)
	favoriteHandler := handler.NewFavouriteHandler(addFavoriteCharacter, removeFavoriteCharacter, fetchFavoriteCharacters)

	// Character endpoints
	http.HandleFunc("/api/character", characterHandler.GetCharacters)

	// User endpoints
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", userHandler.Login)
	http.HandleFunc("/api/logout", userHandler.Logout)

	// Favourites endpoints
	http.HandleFunc("/api/favorites/add", favoriteHandler.AddFavoriteCharacter)
	http.HandleFunc("/api/favorites/remove", favoriteHandler.RemoveFavoriteCharacter)
	http.HandleFunc("/api/favorites", favoriteHandler.FetchFavoriteCharacters)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
