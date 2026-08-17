package handlers

import (
	"net/http"

	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	tagRepo repository.TagRepository
}

func NewTagHandler(tagRepo repository.TagRepository) *TagHandler {
	return &TagHandler{tagRepo: tagRepo}
}

func (h *TagHandler) GetTags(c echo.Context) error {
	tags, err := h.tagRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch tags",
		})
	}

	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	if tagNames == nil {
		tagNames = []string{}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"tags": tagNames,
	})
}
