package handlers

import (
	"net/http"

	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	articleRepo repository.ArticleRepository
}

func NewArticleHandler(articleRepo repository.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{articleRepo: articleRepo}
}

// GetArticles handles GET /api/articles
func (h *ArticleHandler) GetArticles(c echo.Context) error {
	articles, err := h.articleRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch articles",
		})
	}

	articlesDTO := dto.ToArticlesDTO(articles)
	return c.JSON(http.StatusOK, dto.ArticlesResponse{
		Articles:      articlesDTO,
		ArticlesCount: len(articlesDTO),
	})
}

// GetArticleBySlug handles GET /api/articles/:slug
func (h *ArticleHandler) GetArticleBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Slug parameter is required"})
	}

	article, err := h.articleRepo.FindBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Article not found"})
	}

	articleDTO := dto.ToArticleDTO(article)
	return c.JSON(http.StatusOK, dto.ArticleResponse{
		Article: articleDTO,
	})
}
