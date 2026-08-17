package repository

import (
	"realworld-endpoints/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	FindAll() ([]models.Article, error)
	FindBySlug(slug string) (*models.Article, error)
	Create(article *models.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) FindAll() ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Author").Preload("Tags").Find(&articles).Error
	return articles, err
}

func (r *articleRepository) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").Preload("Tags").Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}
