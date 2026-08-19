package dto

import "realworld-endpoints/internal/models"

type UserRegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	User UserRegisterParams `json:"user"`
}

type UserLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	User UserLoginParams `json:"user"`
}

type UserDTO struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserResponse struct {
	User UserDTO `json:"user"`
}

func ToUserDTO(user *models.User, token string) UserDTO {
	return UserDTO{
		Email:    user.Email,
		Token:    token,
		Username: user.Username,
		Bio:      user.Bio,
		Image:    user.Image,
	}
}
