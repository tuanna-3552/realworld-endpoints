package repository

import (
	"errors"
	"realworld-endpoints/internal/models"
)

type MockCommentRepository struct {
	Comments []models.Comment
	NextID   uint
	SlugToArticleID map[string]uint

	ErrFindByArticleSlug error
	ErrFindByID          error
	ErrCreate            error
	ErrDelete            error
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		Comments: []models.Comment{},
		NextID:   1,
		SlugToArticleID: make(map[string]uint),
	}
}

func (m *MockCommentRepository) FindByArticleSlug(slug string) ([]models.Comment, error) {
	if m.ErrFindByArticleSlug != nil {
		return nil, m.ErrFindByArticleSlug
	}
	articleID, ok := m.SlugToArticleID[slug]
	if !ok {
		return []models.Comment{}, nil
	}
	
	var res []models.Comment
	for _, c := range m.Comments {
		if c.ArticleID == articleID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (m *MockCommentRepository) FindByID(id uint) (*models.Comment, error) {
	if m.ErrFindByID != nil {
		return nil, m.ErrFindByID
	}
	for _, c := range m.Comments {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("comment not found")
}

func (m *MockCommentRepository) Create(comment *models.Comment) error {
	if m.ErrCreate != nil {
		return m.ErrCreate
	}
	comment.ID = m.NextID
	m.NextID++
	m.Comments = append(m.Comments, *comment)
	return nil
}

func (m *MockCommentRepository) Delete(id uint) error {
	if m.ErrDelete != nil {
		return m.ErrDelete
	}
	for i, c := range m.Comments {
		if c.ID == id {
			m.Comments = append(m.Comments[:i], m.Comments[i+1:]...)
			return nil
		}
	}
	return errors.New("comment not found")
}
