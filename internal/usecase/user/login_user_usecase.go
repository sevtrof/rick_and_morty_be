package user

import (
	"log"
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepo *repository.UserRepository
}

func NewLoginUserUseCase(userRepo *repository.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{userRepo: userRepo}
}

func (u *LoginUserUseCase) Execute(email, password string) (*model.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding login user by email: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("invalid password: %v", err)
		return nil, err
	}

	return user, nil
}
