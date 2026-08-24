package repository

import (
	"realworld-endpoints/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	FindByArticleSlug(slug string) ([]models.Comment, error)
	FindByID(id uint) (*models.Comment, error)
	Create(comment *models.Comment) error
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) FindByArticleSlug(slug string) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Joins("JOIN articles ON articles.id = comments.article_id").
		Where("articles.slug = ?", slug).
		Preload("Author").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("Author").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Create(comment *models.Comment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return err
	}
	// Preload author for the created comment response
	return r.db.Preload("Author").First(comment, comment.ID).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
