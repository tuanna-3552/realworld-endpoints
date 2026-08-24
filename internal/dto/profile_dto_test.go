package dto

import (
	"testing"
	"realworld-endpoints/internal/models"
)

func TestToProfileDTO_Following(t *testing.T) {
	user := &models.User{
		Username: "jane",
		Bio:      "hi",
		Image:    "image.png",
	}

	result := ToProfileDTO(user, true)

	if result.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, result.Username)
	}
	if result.Bio != user.Bio {
		t.Errorf("expected bio %s, got %s", user.Bio, result.Bio)
	}
	if result.Image != user.Image {
		t.Errorf("expected image %s, got %s", user.Image, result.Image)
	}
	if result.Following != true {
		t.Errorf("expected following true")
	}
}

func TestToProfileDTO_NotFollowing(t *testing.T) {
	user := &models.User{
		Username: "joe",
	}

	result := ToProfileDTO(user, false)

	if result.Following != false {
		t.Errorf("expected following false")
	}
}
