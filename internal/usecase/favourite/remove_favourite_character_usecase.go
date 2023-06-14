package favourite

import (
	"ricknmorty/internal/repository"
)

type RemoveFavouriteCharacter struct {
	userRepo *repository.UserRepository
}

func NewRemoveFavouriteCharacter(userRepo *repository.UserRepository) *RemoveFavouriteCharacter {
	return &RemoveFavouriteCharacter{userRepo: userRepo}
}

func (useCase *RemoveFavouriteCharacter) Execute(userId int, characterId int) error {
	return useCase.userRepo.RemoveFavouriteCharacter(userId, characterId)
}
