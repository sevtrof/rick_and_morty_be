package character

import (
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/repository"
)

type FetchCharactersUsecase struct {
	CharacterRepo *repository.CharacterRepository
}

func NewFetchCharacters(repo *repository.CharacterRepository) *FetchCharactersUsecase {
	return &FetchCharactersUsecase{CharacterRepo: repo}
}

func (uc *FetchCharactersUsecase) Execute(filters map[string]string, page int) (model.CharactersWithInfo, error) {
	return uc.CharacterRepo.FetchCharacters(filters, page)
}
