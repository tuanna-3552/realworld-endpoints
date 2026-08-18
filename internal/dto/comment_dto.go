package dto

import (
	"time"

	"realworld-endpoints/internal/models"
)

type CreateCommentInput struct {
	Body string `json:"body"`
}

type CreateCommentRequest struct {
	Comment CreateCommentInput `json:"comment"`
}

type CommentDTO struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Body      string     `json:"body"`
	Author    ProfileDTO `json:"author"`
}

type CommentResponse struct {
	Comment CommentDTO `json:"comment"`
}

type CommentsResponse struct {
	Comments []CommentDTO `json:"comments"`
}

func ToCommentDTO(comment *models.Comment) CommentDTO {
	return CommentDTO{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Body:      comment.Body,
		Author:    ToProfileDTO(&comment.Author, false),
	}
}

func ToCommentsDTO(comments []models.Comment) []CommentDTO {
	list := make([]CommentDTO, 0, len(comments))
	for _, comment := range comments {
		list = append(list, ToCommentDTO(&comment))
	}
	return list
}
