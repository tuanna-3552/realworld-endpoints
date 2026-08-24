package repository

import (
	"errors"
	"fmt"
	"realworld-endpoints/internal/models"
)

type MockArticleRepository struct {
	Articles  []models.Article
	Favorites map[string]bool
	NextID    uint

	ErrFindAll          error
	ErrFindBySlug       error
	ErrCreate           error
	ErrFindAllWithFilters error
	ErrFindFeed         error
	ErrFavorite         error
	ErrUnfavorite       error
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		Articles:  []models.Article{},
		Favorites: make(map[string]bool),
		NextID:    1,
	}
}

func (m *MockArticleRepository) FindAll() ([]models.Article, error) {
	if m.ErrFindAll != nil {
		return nil, m.ErrFindAll
	}
	return m.Articles, nil
}

func (m *MockArticleRepository) FindBySlug(slug string) (*models.Article, error) {
	if m.ErrFindBySlug != nil {
		return nil, m.ErrFindBySlug
	}
	for _, a := range m.Articles {
		if a.Slug == slug {
			return &a, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *MockArticleRepository) Create(article *models.Article) error {
	if m.ErrCreate != nil {
		return m.ErrCreate
	}
	article.ID = m.NextID
	m.NextID++
	m.Articles = append(m.Articles, *article)
	return nil
}

func (m *MockArticleRepository) FindAllWithFilters(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error) {
	if m.ErrFindAllWithFilters != nil {
		return nil, 0, m.ErrFindAllWithFilters
	}
	total := int64(len(m.Articles))
	
	start := offset
	if start > len(m.Articles) {
		start = len(m.Articles)
	}
	end := start + limit
	if end > len(m.Articles) {
		end = len(m.Articles)
	}
	
	return m.Articles[start:end], total, nil
}

func (m *MockArticleRepository) FindFeed(followedUserIDs []uint, limit, offset int) ([]models.Article, int64, error) {
	if m.ErrFindFeed != nil {
		return nil, 0, m.ErrFindFeed
	}
	
	var filtered []models.Article
	for _, a := range m.Articles {
		isFollowed := false
		for _, id := range followedUserIDs {
			if a.AuthorID == id {
				isFollowed = true
				break
			}
		}
		if isFollowed {
			filtered = append(filtered, a)
		}
	}
	
	total := int64(len(filtered))
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (m *MockArticleRepository) Favorite(userID, articleID uint) error {
	if m.ErrFavorite != nil {
		return m.ErrFavorite
	}
	key := fmt.Sprintf("%d-%d", userID, articleID)
	if !m.Favorites[key] {
		m.Favorites[key] = true
		for i, a := range m.Articles {
			if a.ID == articleID {
				m.Articles[i].FavoritesCount++
				break
			}
		}
	}
	return nil
}

func (m *MockArticleRepository) Unfavorite(userID, articleID uint) error {
	if m.ErrUnfavorite != nil {
		return m.ErrUnfavorite
	}
	key := fmt.Sprintf("%d-%d", userID, articleID)
	if m.Favorites[key] {
		delete(m.Favorites, key)
		for i, a := range m.Articles {
			if a.ID == articleID {
				m.Articles[i].FavoritesCount--
				break
			}
		}
	}
	return nil
}

func (m *MockArticleRepository) IsFavorited(userID, articleID uint) bool {
	key := fmt.Sprintf("%d-%d", userID, articleID)
	return m.Favorites[key]
}
