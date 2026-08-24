package dto

import (
	"testing"
	"realworld-endpoints/internal/models"
)

func TestToCommentDTO(t *testing.T) {
	comment := &models.Comment{
		ID:   1,
		Body: "Great post!",
		Author: models.User{
			Username: "commenter",
		},
	}

	result := ToCommentDTO(comment)

	if result.ID != comment.ID {
		t.Errorf("expected ID %d, got %d", comment.ID, result.ID)
	}
	if result.Body != comment.Body {
		t.Errorf("expected Body %s, got %s", comment.Body, result.Body)
	}
	if result.Author.Username != comment.Author.Username {
		t.Errorf("expected author %s, got %s", comment.Author.Username, result.Author.Username)
	}
	if result.Author.Following != false {
		t.Errorf("expected following false")
	}
}

func TestToCommentDTOWithFollowing(t *testing.T) {
	comment := &models.Comment{
		ID: 2,
		Author: models.User{
			Username: "testuser",
		},
	}

	result := ToCommentDTOWithFollowing(comment, true)

	if result.Author.Following != true {
		t.Errorf("expected following true")
	}
}

func TestToCommentsDTO(t *testing.T) {
	comments := []models.Comment{
		{ID: 1, Body: "first"},
		{ID: 2, Body: "second"},
	}

	dtos := ToCommentsDTO(comments)

	if len(dtos) != 2 {
		t.Fatalf("expected 2 dtos, got %d", len(dtos))
	}
	if dtos[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", dtos[0].ID)
	}
	if dtos[1].ID != 2 {
		t.Errorf("expected ID 2, got %d", dtos[1].ID)
	}
}
