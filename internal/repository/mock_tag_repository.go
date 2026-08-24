package repository

import "realworld-endpoints/internal/models"

type MockTagRepository struct {
	Tags   []models.Tag
	NextID uint

	ErrFindAll error
	ErrCreate  error
}

func NewMockTagRepository() *MockTagRepository {
	return &MockTagRepository{
		Tags:   []models.Tag{},
		NextID: 1,
	}
}

func (m *MockTagRepository) FindAll() ([]models.Tag, error) {
	if m.ErrFindAll != nil {
		return nil, m.ErrFindAll
	}
	return m.Tags, nil
}

func (m *MockTagRepository) Create(tag *models.Tag) error {
	if m.ErrCreate != nil {
		return m.ErrCreate
	}
	tag.ID = m.NextID
	m.NextID++
	m.Tags = append(m.Tags, *tag)
	return nil
}
