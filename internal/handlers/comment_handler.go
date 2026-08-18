package handlers

import (
	"net/http"

	"realworld-endpoints/internal/dto"
	"realworld-endpoints/internal/models"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
}

func NewCommentHandler(
	commentRepo repository.CommentRepository,
	articleRepo repository.ArticleRepository,
	userRepo repository.UserRepository,
) *CommentHandler {
	return &CommentHandler{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

// GetComments handles GET /api/articles/:slug/comments
func (h *CommentHandler) GetComments(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Slug parameter is required"})
	}

	comments, err := h.commentRepo.FindByArticleSlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch comments"})
	}

	commentsDTO := dto.ToCommentsDTO(comments)
	return c.JSON(http.StatusOK, dto.CommentsResponse{Comments: commentsDTO})
}

// CreateComment handles POST /api/articles/:slug/comments
func (h *CommentHandler) CreateComment(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Slug parameter is required"})
	}

	// Find target article
	article, err := h.articleRepo.FindBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Article not found"})
	}

	// Bind request body JSON
	var req dto.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	if req.Comment.Body == "" {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{
				"body": []string{"can't be blank"},
			},
		})
	}

	// Find default author (e.g. johndoe or first user in DB)
	users, err := h.userRepo.FindAll()
	if err != nil || len(users) == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "No user available for comment author"})
	}
	author := users[0]

	comment := models.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  author.ID,
	}

	if err := h.commentRepo.Create(&comment); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create comment"})
	}

	commentDTO := dto.ToCommentDTO(&comment)
	return c.JSON(http.StatusCreated, dto.CommentResponse{Comment: commentDTO})
}
