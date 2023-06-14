package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/usecase/user"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	loginUserUseCase    *user.LoginUserUseCase
	logoutUserUseCase   *user.LogoutUserUseCase
	registerUserUseCase *user.RegisterUserUseCase
}

func NewUserHandler(
	loginUserUseCase *user.LoginUserUseCase,
	logoutUserUseCase *user.LogoutUserUseCase,
	registerUserUseCase *user.RegisterUserUseCase,
) *UserHandler {
	return &UserHandler{loginUserUseCase: loginUserUseCase, logoutUserUseCase: logoutUserUseCase, registerUserUseCase: registerUserUseCase}
}

var jwtKey = []byte("")

func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Registering")
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newUser model.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "error parsing request body", http.StatusBadRequest)
		return
	}

	log.Printf("Registering. Got user: %v", newUser)

	if err := h.generateAvatar(&newUser); err != nil {
		http.Error(w, "Failed to generate avatar", http.StatusInternalServerError)
		return
	}

	log.Printf("Registering. Generated avatar: %v", newUser)

	if err := h.registerUserUseCase.Execute(&newUser); err != nil {
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}

	log.Printf("Registered user: %v", newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": newUser.ID,
		"message": "User successfully registered",
	})
}

func (h *UserHandler) generateAvatar(user *model.User) error {
	log.Printf("Start generating")
	avatarPath := filepath.Join("avatars", fmt.Sprint(user.ID)+".png")

	log.Printf("Generated path: %s", avatarPath)
	resp, err := http.Get("http://localhost:8081/image")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(avatarPath)
	log.Printf("Created os: %v, err: %v", out, err)
	if err != nil {
		log.Printf("failed generating due to: %v", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	user.Avatar = avatarPath

	return nil
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login received: %v", r.Body)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.validateCredentials(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(int(user.ID))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	h.sendResponseWithToken(w, user, token)
}

func (h *UserHandler) validateCredentials(body io.Reader) (*model.User, error) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(body).Decode(&credentials); err != nil {
		return nil, fmt.Errorf("error parsing request body")
	}

	return h.loginUserUseCase.Execute(credentials.Email, credentials.Password)
}

func (h *UserHandler) sendResponseWithToken(w http.ResponseWriter, user *model.User, token string) {
	responseMap := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	responseJSON, err := json.Marshal(responseMap)
	if err != nil {
		http.Error(w, "Error encoding response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {}
