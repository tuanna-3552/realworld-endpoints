package handlers

import (
	"net/http"

	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	userRepo repository.UserRepository
}

func NewProfileHandler(userRepo repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

// GetProfile handles GET /api/profiles/:username
func (h *ProfileHandler) GetProfile(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username parameter is required"})
	}

	user, err := h.userRepo.FindByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Profile not found"})
	}

	profileDTO := dto.ToProfileDTO(user, false)
	return c.JSON(http.StatusOK, dto.ProfileResponse{Profile: profileDTO})
}
