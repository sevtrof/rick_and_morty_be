package user

import (
	"errors"
	"log"
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserUseCase struct {
	userRepo *repository.UserRepository
}

func NewRegisterUserUseCase(userRepo *repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: userRepo}
}

func (u *RegisterUserUseCase) Execute(user *model.User) error {

	log.Printf("Registering in usecase: %v", *user)
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		if err.Error() != "user not found" {
			log.Printf("error finding user by email: %v", err)
			return err
		}
	}

	if existingUser != nil {
		return errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error generating hashed password: %v", err)
		return err
	}

	user.Password = string(hashedPassword)

	return u.userRepo.Save(user)
}
