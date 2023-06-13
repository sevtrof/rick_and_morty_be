package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/usecase/user"
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

	log.Printf("Registering user: %v", newUser)

	if err := h.registerUserUseCase.Execute(&newUser); err != nil {
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}

	log.Printf("Registered user: %v", newUser)

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "User registered successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	user, err := h.loginUserUseCase.Execute(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	token := "dummy_token"

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
