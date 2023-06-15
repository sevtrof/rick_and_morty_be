package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ricknmorty/internal/usecase/news"
)

type NewsHandler struct {
	fetchNewsUseCase *news.FetchNewsUsecase
}

func NewNewsHandler(fetchNewsUseCase *news.FetchNewsUsecase) *NewsHandler {
	return &NewsHandler{fetchNewsUseCase: fetchNewsUseCase}
}

func (h *NewsHandler) FetchNews(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	newsItems, err := h.fetchNewsUseCase.FetchNews(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsItems)
}
