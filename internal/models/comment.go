package models

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	ArticleID uint      `gorm:"not null;index" json:"articleId"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"-"`
	AuthorID  uint      `gorm:"not null" json:"authorId"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
