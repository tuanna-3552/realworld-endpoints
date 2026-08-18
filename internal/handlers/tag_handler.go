package handlers

import (
	"net/http"

	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	tagRepo repository.TagRepository
}

func NewTagHandler(tagRepo repository.TagRepository) *TagHandler {
	return &TagHandler{tagRepo: tagRepo}
}

// GetTags handles GET /api/tags
func (h *TagHandler) GetTags(c echo.Context) error {
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

	return c.JSON(http.StatusOK, dto.TagsResponse{
		Tags: tagNames,
	})
}
