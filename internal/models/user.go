package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:255;not null;uniqueIndex" json:"username"`
	Email     string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Bio       string    `gorm:"type:text" json:"bio"`
	Image     string    `gorm:"size:512" json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
