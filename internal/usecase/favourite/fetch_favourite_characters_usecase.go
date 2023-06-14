package favourite

import (
	"ricknmorty/internal/repository"
)

type FetchFavouriteCharacters struct {
	userRepo *repository.UserRepository
}

func NewFetchFavouriteCharacters(userRepo *repository.UserRepository) *FetchFavouriteCharacters {
	return &FetchFavouriteCharacters{userRepo: userRepo}
}

func (useCase *FetchFavouriteCharacters) Execute(userId int) ([]int, error) {
	user, err := useCase.userRepo.FindById(userId)
	if err != nil {
		return nil, err
	}
	return user.FavouriteCharacters, nil
}
