package repository

import (
	"realworld-endpoints/internal/models"

	"gorm.io/gorm"
)

type TagRepository interface {
	FindAll() ([]models.Tag, error)
	Create(tag *models.Tag) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}
