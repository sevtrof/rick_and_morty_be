package news

import (
	"ricknmorty/internal/domain/model"
	"ricknmorty/internal/repository"
)

type FetchNewsUsecase struct {
	newsRepo *repository.NewsRepository
}

func NewFetchNewsUseCase(newsRepo *repository.NewsRepository) *FetchNewsUsecase {
	return &FetchNewsUsecase{newsRepo: newsRepo}
}

func (uc *FetchNewsUsecase) FetchNews(page int) ([]*model.NewsItem, error) {
	return uc.newsRepo.FetchNews(page)
}
