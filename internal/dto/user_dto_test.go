package dto

import (
	"testing"
	"realworld-endpoints/internal/models"
)

func TestToUserDTO(t *testing.T) {
	user := &models.User{
		Username: "johndoe",
		Email:    "john@example.com",
		Bio:      "I am a developer",
		Image:    "http://image.com/img.png",
	}
	token := "jwt-token"

	result := ToUserDTO(user, token)

	if result.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, result.Username)
	}
	if result.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, result.Email)
	}
	if result.Bio != user.Bio {
		t.Errorf("expected bio %s, got %s", user.Bio, result.Bio)
	}
	if result.Image != user.Image {
		t.Errorf("expected image %s, got %s", user.Image, result.Image)
	}
	if result.Token != token {
		t.Errorf("expected token %s, got %s", token, result.Token)
	}
}
