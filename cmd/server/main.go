package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"ricknmorty/internal/domain/service"
	"ricknmorty/internal/handler"
	"ricknmorty/internal/repository"
	"ricknmorty/internal/usecase/character"
	"ricknmorty/internal/usecase/favourite"
	"ricknmorty/internal/usecase/news"
	"ricknmorty/internal/usecase/user"

	"github.com/joho/godotenv"
)

const (
	avatarPath string = "/app/avatars"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	tokenSecretKey := os.Getenv("SECRET_KEY")
	connStr := os.Getenv("DB_CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// News
	newsRepo := repository.NewNewsRepository()
	fetchNewsUseCase := news.NewFetchNewsUseCase(newsRepo)
	newsHandler := handler.NewNewsHandler(fetchNewsUseCase)

	// Service
	avatarGen := service.NewAvatarService()
	tokenGen := service.NewTokenService(tokenSecretKey)
	tokenParser := service.NewJwtTokenParser(tokenSecretKey)

	// Character
	characterRepo := repository.NewCharacterRepository(db)
	fetchCharacters := character.NewFetchCharacters(characterRepo)
	characterHandler := handler.NewCharacterHandler(fetchCharacters)

	// User
	userRepo := repository.NewUserRepository(db, avatarGen)
	loginUserUseCase := user.NewLoginUserUseCase(userRepo)
	logoutUserUseCase := user.NewLogoutUserUseCase()
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(loginUserUseCase, logoutUserUseCase, registerUserUseCase, tokenGen)

	// Favourite
	addFavoriteCharacter := favourite.NewAddFavouriteCharacter(userRepo)
	removeFavoriteCharacter := favourite.NewRemoveFavouriteCharacter(userRepo)
	fetchFavoriteCharacters := favourite.NewFetchFavouriteCharacters(userRepo)
	favoriteHandler := handler.NewFavouriteHandler(addFavoriteCharacter, removeFavoriteCharacter, fetchFavoriteCharacters, tokenParser)

	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir(avatarPath))))

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

	// News endpoints
	http.HandleFunc("/api/news", newsHandler.FetchNews)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
