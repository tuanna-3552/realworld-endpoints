package dto

import "realworld-endpoints/internal/models"

type ProfileDTO struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ProfileResponse struct {
	Profile ProfileDTO `json:"profile"`
}

func ToProfileDTO(user *models.User, following bool) ProfileDTO {
	return ProfileDTO{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: following,
	}
}
