package repository

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"ricknmorty/internal/domain/model"

	"github.com/joho/godotenv"
)

type NewsRepository struct{}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{}
}

func (r *NewsRepository) FetchNews(page int) ([]*model.NewsItem, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	const pageSize = 20
	resp, err := http.Get(os.Getenv("NEWS_GEN_URL"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var newsResponse map[string][]string
	json.Unmarshal(body, &newsResponse)

	var newsItems []*model.NewsItem
	for _, content := range newsResponse["news"] {
		newsItems = append(newsItems, &model.NewsItem{Content: content})
	}

	start := (page - 1) * pageSize
	if start >= len(newsItems) {
		return []*model.NewsItem{}, nil
	}

	end := start + pageSize
	if end > len(newsItems) {
		end = len(newsItems)
	}

	return newsItems[start:end], nil
}
