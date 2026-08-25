package handlers

import (
	"net/http"
	"strconv"

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
// @Summary      Get comments for an article
// @Description  Get all comments for an article identified by slug
// @Tags         Comments
// @Produce      json
// @Param        slug  path  string  true  "Article slug"
// @Success      200  {object}  dto.CommentsResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /articles/{slug}/comments [get]
func (h *CommentHandler) GetComments(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Slug parameter is required"}}})
	}

	comments, err := h.commentRepo.FindByArticleSlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"Failed to fetch comments"}}})
	}

	var currentUserID uint
	if currentUserIDVal := c.Get("user_id"); currentUserIDVal != nil {
		if id, ok := currentUserIDVal.(uint); ok {
			currentUserID = id
		}
	}

	authorIDs := make([]uint, len(comments))
	for i, c := range comments {
		authorIDs[i] = c.AuthorID
	}
	folMap, _ := h.userRepo.BatchIsFollowing(currentUserID, authorIDs)

	commentsDTO := make([]dto.CommentDTO, 0, len(comments))
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, dto.ToCommentDTOWithFollowing(
			&comment, folMap[comment.AuthorID],
		))
	}

	return c.JSON(http.StatusOK, dto.CommentsResponse{Comments: commentsDTO})
}

// CreateComment handles POST /api/articles/:slug/comments
// @Summary      Add comment to an article
// @Description  Create a new comment on an article
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        slug     path  string                    true  "Article slug"
// @Param        comment  body  dto.CreateCommentRequest  true  "Comment payload"
// @Success      201  {object}  dto.CommentResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /articles/{slug}/comments [post]
func (h *CommentHandler) CreateComment(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Slug parameter is required"}}})
	}

	article, err := h.articleRepo.FindBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"article": []string{"not found"}}})
	}

	var req dto.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Invalid request payload"}}})
	}

	if req.Comment.Body == "" {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"errors": echo.Map{
				"body": []string{"can't be blank"},
			},
		})
	}

	var authorID uint
	if currentUserIDVal := c.Get("user_id"); currentUserIDVal != nil {
		authorID = currentUserIDVal.(uint)
	} else {
		users, err := h.userRepo.FindAll()
		if err != nil || len(users) == 0 {
			return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"No user available for comment author"}}})
		}
		authorID = users[0].ID
	}

	comment := models.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  authorID,
	}

	if err := h.commentRepo.Create(&comment); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"Failed to create comment"}}})
	}

	commentDTO := dto.ToCommentDTO(&comment)
	return c.JSON(http.StatusCreated, dto.CommentResponse{Comment: commentDTO})
}

// DeleteComment handles DELETE /api/articles/:slug/comments/:id (Ownership check)
// @Summary      Delete a comment
// @Description  Delete a comment (only the author of the comment can delete it)
// @Tags         Comments
// @Produce      json
// @Security     ApiKeyAuth
// @Param        slug  path  string  true  "Article slug"
// @Param        id    path  int     true  "Comment ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /articles/{slug}/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Slug parameter is required"}}})
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || commentID <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": echo.Map{"body": []string{"Invalid comment ID"}}})
	}

	currentUserIDVal := c.Get("user_id")
	if currentUserIDVal == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"errors": echo.Map{"body": []string{"unauthorized"}}})
	}
	currentUserID := currentUserIDVal.(uint)

	comment, err := h.commentRepo.FindByID(uint(commentID))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"errors": echo.Map{"comment": []string{"not found"}}})
	}

	// Ownership check: only author of the comment can delete it
	if comment.AuthorID != currentUserID {
		return c.JSON(http.StatusForbidden, echo.Map{
			"errors": echo.Map{"comment": []string{"you are not authorized to delete this comment"}},
		})
	}

	if err := h.commentRepo.Delete(uint(commentID)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"errors": echo.Map{"body": []string{"Failed to delete comment"}}})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "comment deleted successfully"})
}
