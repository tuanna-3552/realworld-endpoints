package handlers

import (
	"net/http"

	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch users",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"users": users,
		"count": len(users),
	})
}
