package handlers

import (
	"net/http"

	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	articleRepo repository.ArticleRepository
}

func NewArticleHandler(articleRepo repository.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{articleRepo: articleRepo}
}

func (h *ArticleHandler) GetArticles(c echo.Context) error {
	articles, err := h.articleRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch articles",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"articles":      articles,
		"articlesCount": len(articles),
	})
}
