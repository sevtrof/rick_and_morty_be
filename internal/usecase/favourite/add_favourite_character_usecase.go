package favourite

import (
	"ricknmorty/internal/repository"
)

type AddFavouriteCharacter struct {
	userRepo *repository.UserRepository
}

func NewAddFavouriteCharacter(userRepo *repository.UserRepository) *AddFavouriteCharacter {
	return &AddFavouriteCharacter{userRepo: userRepo}
}

func (useCase *AddFavouriteCharacter) Execute(userId int, characterId int) error {
	return useCase.userRepo.AddFavouriteCharacter(userId, characterId)
}
