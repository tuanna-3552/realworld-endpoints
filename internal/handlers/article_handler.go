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
// @Summary      List articles
// @Description  List articles with optional filtering by tag, author, favorited user, and pagination
// @Tags         Articles
// @Produce      json
// @Param        tag        query   string  false  "Filter by tag"
// @Param        author     query   string  false  "Filter by author username"
// @Param        favorited  query   string  false  "Filter by favorited username"
// @Param        limit      query   int     false  "Limit number of articles (default 20)"
// @Param        offset     query   int     false  "Offset/skip number of articles (default 0)"
// @Success      200  {object}  dto.ArticlesResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /articles [get]
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

	articleIDs := make([]uint, len(articles))
	authorIDs := make([]uint, len(articles))
	for i, a := range articles {
		articleIDs[i] = a.ID
		authorIDs[i] = a.AuthorID
	}
	favMap, _ := h.articleRepo.BatchIsFavorited(currentUserID, articleIDs)
	folMap, _ := h.userRepo.BatchIsFollowing(currentUserID, authorIDs)

	articlesDTO := make([]dto.ArticleDTO, 0, len(articles))
	for _, article := range articles {
		articlesDTO = append(articlesDTO, dto.ToArticleDTOWithStatus(
			&article, favMap[article.ID], folMap[article.AuthorID],
		))
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
// @Summary      Feed articles
// @Description  Get articles from users that the current user follows
// @Tags         Articles
// @Produce      json
// @Security     ApiKeyAuth
// @Param        limit   query   int  false  "Limit number of articles (default 20)"
// @Param        offset  query   int  false  "Offset/skip number of articles (default 0)"
// @Success      200  {object}  dto.ArticlesResponse
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /articles/feed [get]
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

	articleIDs := make([]uint, len(articles))
	for i, a := range articles {
		articleIDs[i] = a.ID
	}
	favMap, _ := h.articleRepo.BatchIsFavorited(currentUserID, articleIDs)

	articlesDTO := make([]dto.ArticleDTO, 0, len(articles))
	for _, article := range articles {
		articlesDTO = append(articlesDTO, dto.ToArticleDTOWithStatus(
			&article, favMap[article.ID], true,
		))
	}

	return c.JSON(http.StatusOK, dto.ArticlesResponse{
		Articles:      articlesDTO,
		ArticlesCount: int(count),
	})
}

// GetArticleBySlug handles GET /api/articles/:slug
// @Summary      Get article by slug
// @Description  Get a single article by its slug
// @Tags         Articles
// @Produce      json
// @Param        slug  path  string  true  "Article slug"
// @Success      200  {object}  dto.ArticleResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /articles/{slug} [get]
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
// @Summary      Favorite an article
// @Description  Mark an article as favorite for the current user
// @Tags         Favorites
// @Produce      json
// @Security     ApiKeyAuth
// @Param        slug  path  string  true  "Article slug"
// @Success      200  {object}  dto.ArticleResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /articles/{slug}/favorite [post]
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
// @Summary      Unfavorite an article
// @Description  Remove an article from the current user's favorites
// @Tags         Favorites
// @Produce      json
// @Security     ApiKeyAuth
// @Param        slug  path  string  true  "Article slug"
// @Success      200  {object}  dto.ArticleResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /articles/{slug}/favorite [delete]
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
