package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/domain/service"
	"ricknmorty/internal/usecase/user"
)

type UserHandler struct {
	loginUserUseCase    *user.LoginUserUseCase
	logoutUserUseCase   *user.LogoutUserUseCase
	registerUserUseCase *user.RegisterUserUseCase
	tokenService        *service.TokenService
}

func NewUserHandler(
	loginUserUseCase *user.LoginUserUseCase,
	logoutUserUseCase *user.LogoutUserUseCase,
	registerUserUseCase *user.RegisterUserUseCase,
	tokenService *service.TokenService,
) *UserHandler {
	return &UserHandler{
		loginUserUseCase:    loginUserUseCase,
		logoutUserUseCase:   logoutUserUseCase,
		registerUserUseCase: registerUserUseCase,
		tokenService:        tokenService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newUser model.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	err := h.registerUserUseCase.Execute(&newUser)
	if err != nil {
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": newUser.ID,
		"message": "User successfully registered",
	})
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

	token, err := h.tokenService.GenerateToken(int(user.ID))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	baseAvatarURL := "http://localhost:8080/" + user.Avatar
	user.Avatar = baseAvatarURL

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
