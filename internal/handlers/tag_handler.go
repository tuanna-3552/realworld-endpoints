package handlers

import (
	"encoding/json"
	"net/http"

	"realworld-endpoints/internal/cache"
	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	tagRepo repository.TagRepository
	cache   cache.CacheService
}

func NewTagHandler(tagRepo repository.TagRepository, cacheService cache.CacheService) *TagHandler {
	return &TagHandler{
		tagRepo: tagRepo,
		cache:   cacheService,
	}
}

// GetTags handles GET /api/tags
// @Summary      Get tags
// @Description  Get a list of all tags used across articles
// @Tags         Tags
// @Produce      json
// @Success      200  {object}  dto.TagsResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /tags [get]
func (h *TagHandler) GetTags(c echo.Context) error {
	ctx := c.Request().Context()
	
	if h.cache != nil {
		if cached, err := h.cache.Get(ctx, cache.TagsCacheKey); err == nil {
			var resp dto.TagsResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				return c.JSON(http.StatusOK, resp)
			}
		}
	}

	tags, err := h.tagRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch tags",
		})
	}

	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	resp := dto.TagsResponse{
		Tags: tagNames,
	}

	if h.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			_ = h.cache.Set(ctx, cache.TagsCacheKey, string(data), cache.TagsCacheTTL)
		}
	}

	return c.JSON(http.StatusOK, resp)
}
