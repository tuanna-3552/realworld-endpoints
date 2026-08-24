package models

import (
	"time"
)

type Article struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Slug           string    `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Title          string    `gorm:"size:255;not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	Body           string    `gorm:"type:text" json:"body"`
	FavoritesCount int       `gorm:"default:0" json:"favoritesCount"`
	AuthorID       uint      `gorm:"not null" json:"authorId"`
	Author         User      `gorm:"foreignKey:AuthorID" json:"author"`
	Tags           []Tag     `gorm:"many2many:article_tags;" json:"tags"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ArticleFavorite struct {
	UserID    uint      `gorm:"primaryKey;index" json:"userId"`
	ArticleID uint      `gorm:"primaryKey;index" json:"articleId"`
	CreatedAt time.Time `json:"createdAt"`
}
