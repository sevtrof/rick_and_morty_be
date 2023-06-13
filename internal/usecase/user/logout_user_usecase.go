package user

import (
	"ricknmorty/internal/domain/model"
)

type LogoutUserUseCase struct{}

func NewLogoutUserUseCase() *LogoutUserUseCase {
	return &LogoutUserUseCase{}
}

func (u *LogoutUserUseCase) Execute(user *model.User) error {
	return nil
}
