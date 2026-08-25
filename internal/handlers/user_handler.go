package handlers

import (
	"net/http"
	"strings"

	"realworld-endpoints/internal/auth"
	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/models"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewUserHandler(userRepo repository.UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// GetUsers handles GET /api/users (List all users)
// @Summary      List all users
// @Description  Returns a list of all registered users
// @Tags         Users
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users [get]
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

// Register handles POST /api/users (User Registration)
// @Summary      Register a new user
// @Description  Create a new user account and return JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserRegisterRequest  true  "User registration payload"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      422   {object}  map[string]interface{}
// @Router       /users [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req dto.UserRegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": echo.Map{"body": []string{"invalid request payload"}},
		})
	}

	if req.User.Username == "" || req.User.Email == "" || req.User.Password == "" {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{"body": []string{"username, email, and password are required"}},
		})
	}

	// Check if username already exists
	if existing, _ := h.userRepo.FindByUsername(req.User.Username); existing != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{"username": []string{"has already been taken"}},
		})
	}

	// Check if email already exists
	if existing, _ := h.userRepo.FindByEmail(req.User.Email); existing != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{"email": []string{"has already been taken"}},
		})
	}

	hashedPassword, err := auth.HashPassword(req.User.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"failed to hash password"}},
		})
	}

	user := models.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: hashedPassword,
	}

	if err := h.userRepo.Create(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"failed to create user"}},
		})
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"failed to generate token"}},
		})
	}

	return c.JSON(http.StatusCreated, dto.UserResponse{
		User: dto.ToUserDTO(&user, token),
	})
}

// Login handles POST /api/users/login (User Authentication)
// @Summary      User login
// @Description  Authenticate a user and return JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserLoginRequest  true  "User login credentials"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      422   {object}  map[string]interface{}
// @Router       /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req dto.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": echo.Map{"body": []string{"invalid request payload"}},
		})
	}

	if req.User.Email == "" || req.User.Password == "" {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{"body": []string{"email and password are required"}},
		})
	}

	user, err := h.userRepo.FindByEmail(req.User.Email)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{"email or password": []string{"is invalid"}},
		})
	}

	if !auth.CheckPasswordHash(req.User.Password, user.Password) {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{"email or password": []string{"is invalid"}},
		})
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": echo.Map{"body": []string{"failed to generate token"}},
		})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		User: dto.ToUserDTO(user, token),
	})
}

// GetCurrentUser handles GET /api/user (Get current user by token)
// @Summary      Get current user
// @Description  Get the currently authenticated user profile
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /user [get]
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"errors": echo.Map{"body": []string{"unauthorized"}},
		})
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"errors": echo.Map{"body": []string{"invalid user session"}},
		})
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"errors": echo.Map{"user": []string{"not found"}},
		})
	}

	authHeader := c.Request().Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	token := ""
	if len(parts) == 2 {
		token = parts[1]
	} else {
		token, _ = auth.GenerateToken(user.ID, user.Email, h.jwtSecret)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		User: dto.ToUserDTO(user, token),
	})
}
