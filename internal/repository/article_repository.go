package repository

import (
	"time"

	"realworld-endpoints/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	FindAll() ([]models.Article, error)
	FindBySlug(slug string) (*models.Article, error)
	Create(article *models.Article) error

	// Filtering & Pagination
	FindAllWithFilters(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error)
	FindFeed(followedUserIDs []uint, limit, offset int) ([]models.Article, int64, error)

	// Article Favorites System
	Favorite(userID, articleID uint) error
	Unfavorite(userID, articleID uint) error
	IsFavorited(userID, articleID uint) bool
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

func (r *articleRepository) FindAllWithFilters(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error) {
	query := r.db.Model(&models.Article{})

	if tag != "" {
		query = query.Where("id IN (SELECT article_id FROM article_tags JOIN tags ON tags.id = article_tags.tag_id WHERE tags.name = ?)", tag)
	}

	if author != "" {
		query = query.Where("author_id IN (SELECT id FROM users WHERE username = ?)", author)
	}

	if favorited != "" {
		query = query.Where("id IN (SELECT article_id FROM article_favorites JOIN users ON users.id = article_favorites.user_id WHERE users.username = ?)", favorited)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}

	var articles []models.Article
	err := query.Preload("Author").Preload("Tags").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error

	return articles, count, err
}

func (r *articleRepository) FindFeed(followedUserIDs []uint, limit, offset int) ([]models.Article, int64, error) {
	if len(followedUserIDs) == 0 {
		return []models.Article{}, 0, nil
	}

	query := r.db.Model(&models.Article{}).Where("author_id IN ?", followedUserIDs)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}

	var articles []models.Article
	err := query.Preload("Author").Preload("Tags").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error

	return articles, count, err
}

func (r *articleRepository) Favorite(userID, articleID uint) error {
	fav := models.ArticleFavorite{
		UserID:    userID,
		ArticleID: articleID,
		CreatedAt: time.Now(),
	}

	err := r.db.Where(models.ArticleFavorite{UserID: userID, ArticleID: articleID}).
		FirstOrCreate(&fav).Error
	if err != nil {
		return err
	}

	// Update favorites count in article
	var favCount int64
	r.db.Model(&models.ArticleFavorite{}).Where("article_id = ?", articleID).Count(&favCount)
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).Update("favorites_count", favCount).Error
}

func (r *articleRepository) Unfavorite(userID, articleID uint) error {
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&models.ArticleFavorite{}).Error
	if err != nil {
		return err
	}

	// Update favorites count in article
	var favCount int64
	r.db.Model(&models.ArticleFavorite{}).Where("article_id = ?", articleID).Count(&favCount)
	return r.db.Model(&models.Article{}).Where("id = ?", articleID).Update("favorites_count", favCount).Error
}

func (r *articleRepository) IsFavorited(userID, articleID uint) bool {
	if userID == 0 || articleID == 0 {
		return false
	}
	var count int64
	r.db.Model(&models.ArticleFavorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count)
	return count > 0
}
