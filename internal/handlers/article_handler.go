package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"realworld-endpoints/internal/cache"
	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
	cache       cache.CacheService
}

func NewArticleHandler(articleRepo repository.ArticleRepository, userRepo repository.UserRepository, cacheService cache.CacheService) *ArticleHandler {
	return &ArticleHandler{
		articleRepo: articleRepo,
		userRepo:    userRepo,
		cache:       cacheService,
	}
}

// GetArticles handles GET /api/articles (Filter & Pagination)
func (h *ArticleHandler) GetArticles(c echo.Context) error {
	tag := c.QueryParam("tag")
	author := c.QueryParam("author")
	favorited := c.QueryParam("favorited")

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var currentUserID uint
	if currentUserIDVal := c.Get("user_id"); currentUserIDVal != nil {
		if id, ok := currentUserIDVal.(uint); ok {
			currentUserID = id
		}
	}

	ctx := c.Request().Context()
	cacheKey := cache.ArticlesCacheKey(tag, author, favorited, limit, offset)

	if currentUserID == 0 && h.cache != nil {
		if cached, err := h.cache.Get(ctx, cacheKey); err == nil {
			var resp dto.ArticlesResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				return c.JSON(http.StatusOK, resp)
			}
		}
	}

	articles, count, err := h.articleRepo.FindAllWithFilters(tag, author, favorited, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"Failed to fetch articles"}},
		})
	}

	articlesDTO := make([]dto.ArticleDTO, 0, len(articles))
	for _, article := range articles {
		isFav := h.articleRepo.IsFavorited(currentUserID, article.ID)
		isFol := h.userRepo.IsFollowing(currentUserID, article.AuthorID)
		articlesDTO = append(articlesDTO, dto.ToArticleDTOWithStatus(&article, isFav, isFol))
	}

	resp := dto.ArticlesResponse{
		Articles:      articlesDTO,
		ArticlesCount: int(count),
	}

	if currentUserID == 0 && h.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			_ = h.cache.Set(ctx, cacheKey, string(data), cache.ArticlesCacheTTL)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// GetFeed handles GET /api/articles/feed (User Feed for followed authors)
func (h *ArticleHandler) GetFeed(c echo.Context) error {
	currentUserIDVal := c.Get("user_id")
	if currentUserIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"errors": echo.Map{"body": []string{"unauthorized"}},
		})
	}
	currentUserID := currentUserIDVal.(uint)

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	followedIDs, err := h.userRepo.GetFollowedUserIDs(currentUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"Failed to fetch followed users"}},
		})
	}

	articles, count, err := h.articleRepo.FindFeed(followedIDs, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"Failed to fetch feed articles"}},
		})
	}

	articlesDTO := make([]dto.ArticleDTO, 0, len(articles))
	for _, article := range articles {
		isFav := h.articleRepo.IsFavorited(currentUserID, article.ID)
		articlesDTO = append(articlesDTO, dto.ToArticleDTOWithStatus(&article, isFav, true))
	}

	return c.JSON(http.StatusOK, dto.ArticlesResponse{
		Articles:      articlesDTO,
		ArticlesCount: int(count),
	})
}

// GetArticleBySlug handles GET /api/articles/:slug
func (h *ArticleHandler) GetArticleBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Slug parameter is required"}}})
	}

	var currentUserID uint
	if currentUserIDVal := c.Get("user_id"); currentUserIDVal != nil {
		if id, ok := currentUserIDVal.(uint); ok {
			currentUserID = id
		}
	}

	ctx := c.Request().Context()
	cacheKey := cache.ArticleCacheKey(slug)

	if currentUserID == 0 && h.cache != nil {
		if cached, err := h.cache.Get(ctx, cacheKey); err == nil {
			var resp dto.ArticleResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				return c.JSON(http.StatusOK, resp)
			}
		}
	}

	article, err := h.articleRepo.FindBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"article": []string{"not found"}}})
	}

	isFav := h.articleRepo.IsFavorited(currentUserID, article.ID)
	isFol := h.userRepo.IsFollowing(currentUserID, article.AuthorID)

	articleDTO := dto.ToArticleDTOWithStatus(article, isFav, isFol)
	resp := dto.ArticleResponse{
		Article: articleDTO,
	}

	if currentUserID == 0 && h.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			_ = h.cache.Set(ctx, cacheKey, string(data), cache.ArticleCacheTTL)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// FavoriteArticle handles POST /api/articles/:slug/favorite
func (h *ArticleHandler) FavoriteArticle(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Slug parameter is required"}}})
	}

	currentUserIDVal := c.Get("user_id")
	if currentUserIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"errors": echo.Map{"body": []string{"unauthorized"}}})
	}
	currentUserID := currentUserIDVal.(uint)

	article, err := h.articleRepo.FindBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"article": []string{"not found"}}})
	}

	if err := h.articleRepo.Favorite(currentUserID, article.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"failed to favorite article"}}})
	}

	// Invalidate cache
	if h.cache != nil {
		ctx := c.Request().Context()
		_ = h.cache.Delete(ctx, cache.ArticleCacheKey(slug))
		_ = h.cache.DeleteByPattern(ctx, cache.ArticlesCachePrefix+"*")
	}

	// Reload article
	article, _ = h.articleRepo.FindBySlug(slug)
	isFol := h.userRepo.IsFollowing(currentUserID, article.AuthorID)

	return c.JSON(http.StatusOK, dto.ArticleResponse{
		Article: dto.ToArticleDTOWithStatus(article, true, isFol),
	})
}

// UnfavoriteArticle handles DELETE /api/articles/:slug/favorite
func (h *ArticleHandler) UnfavoriteArticle(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Slug parameter is required"}}})
	}

	currentUserIDVal := c.Get("user_id")
	if currentUserIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"errors": echo.Map{"body": []string{"unauthorized"}}})
	}
	currentUserID := currentUserIDVal.(uint)

	article, err := h.articleRepo.FindBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"article": []string{"not found"}}})
	}

	if err := h.articleRepo.Unfavorite(currentUserID, article.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"failed to unfavorite article"}}})
	}

	// Invalidate cache
	if h.cache != nil {
		ctx := c.Request().Context()
		_ = h.cache.Delete(ctx, cache.ArticleCacheKey(slug))
		_ = h.cache.DeleteByPattern(ctx, cache.ArticlesCachePrefix+"*")
	}

	// Reload article
	article, _ = h.articleRepo.FindBySlug(slug)
	isFol := h.userRepo.IsFollowing(currentUserID, article.AuthorID)

	return c.JSON(http.StatusOK, dto.ArticleResponse{
		Article: dto.ToArticleDTOWithStatus(article, false, isFol),
	})
}
