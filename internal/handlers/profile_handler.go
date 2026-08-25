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
// @Summary      Get user profile
// @Description  Get a user's public profile by username
// @Tags         Profiles
// @Produce      json
// @Param        username  path  string  true  "Username"
// @Success      200  {object}  dto.ProfileResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /profiles/{username} [get]
func (h *ProfileHandler) GetProfile(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Username parameter is required"}}})
	}

	user, err := h.userRepo.FindByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"profile": []string{"not found"}}})
	}

	following := false
	if currentUserIDVal := c.Get("user_id"); currentUserIDVal != nil {
		if currentUserID, ok := currentUserIDVal.(uint); ok {
			following = h.userRepo.IsFollowing(currentUserID, user.ID)
		}
	}

	profileDTO := dto.ToProfileDTO(user, following)
	return c.JSON(http.StatusOK, dto.ProfileResponse{Profile: profileDTO})
}

// FollowUser handles POST /api/profiles/:username/follow
// @Summary      Follow a user
// @Description  Follow a user by username
// @Tags         Profiles
// @Produce      json
// @Security     ApiKeyAuth
// @Param        username  path  string  true  "Username to follow"
// @Success      200  {object}  dto.ProfileResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /profiles/{username}/follow [post]
func (h *ProfileHandler) FollowUser(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Username parameter is required"}}})
	}

	currentUserIDVal := c.Get("user_id")
	if currentUserIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"errors": echo.Map{"body": []string{"unauthorized"}}})
	}
	currentUserID := currentUserIDVal.(uint)

	targetUser, err := h.userRepo.FindByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"profile": []string{"not found"}}})
	}

	if err := h.userRepo.Follow(currentUserID, targetUser.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"failed to follow user"}}})
	}

	profileDTO := dto.ToProfileDTO(targetUser, true)
	return c.JSON(http.StatusOK, dto.ProfileResponse{Profile: profileDTO})
}

// UnfollowUser handles DELETE /api/profiles/:username/follow
// @Summary      Unfollow a user
// @Description  Unfollow a user by username
// @Tags         Profiles
// @Produce      json
// @Security     ApiKeyAuth
// @Param        username  path  string  true  "Username to unfollow"
// @Success      200  {object}  dto.ProfileResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /profiles/{username}/follow [delete]
func (h *ProfileHandler) UnfollowUser(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Username parameter is required"}}})
	}

	currentUserIDVal := c.Get("user_id")
	if currentUserIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"errors": echo.Map{"body": []string{"unauthorized"}}})
	}
	currentUserID := currentUserIDVal.(uint)

	targetUser, err := h.userRepo.FindByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"profile": []string{"not found"}}})
	}

	if err := h.userRepo.Unfollow(currentUserID, targetUser.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"failed to unfollow user"}}})
	}

	profileDTO := dto.ToProfileDTO(targetUser, false)
	return c.JSON(http.StatusOK, dto.ProfileResponse{Profile: profileDTO})
}
