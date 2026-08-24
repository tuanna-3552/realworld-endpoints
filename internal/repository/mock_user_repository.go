package repository

import (
	"errors"
	"fmt"
	"realworld-endpoints/internal/models"
)

type MockUserRepository struct {
	Users []models.User
	Follows map[string]bool
	NextID uint

	ErrFindAll error
	ErrFindByID error
	ErrFindByUsername error
	ErrFindByEmail error
	ErrCreate error
	ErrFollow error
	ErrUnfollow error
	ErrGetFollowedUserIDs error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: []models.User{},
		Follows: make(map[string]bool),
		NextID: 1,
	}
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	if m.ErrFindAll != nil {
		return nil, m.ErrFindAll
	}
	return m.Users, nil
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	if m.ErrFindByID != nil {
		return nil, m.ErrFindByID
	}
	for _, u := range m.Users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	if m.ErrFindByUsername != nil {
		return nil, m.ErrFindByUsername
	}
	for _, u := range m.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	if m.ErrFindByEmail != nil {
		return nil, m.ErrFindByEmail
	}
	for _, u := range m.Users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.ErrCreate != nil {
		return m.ErrCreate
	}
	user.ID = m.NextID
	m.NextID++
	m.Users = append(m.Users, *user)
	return nil
}

func (m *MockUserRepository) Follow(followerID, followedID uint) error {
	if m.ErrFollow != nil {
		return m.ErrFollow
	}
	key := fmt.Sprintf("%d-%d", followerID, followedID)
	m.Follows[key] = true
	return nil
}

func (m *MockUserRepository) Unfollow(followerID, followedID uint) error {
	if m.ErrUnfollow != nil {
		return m.ErrUnfollow
	}
	key := fmt.Sprintf("%d-%d", followerID, followedID)
	delete(m.Follows, key)
	return nil
}

func (m *MockUserRepository) IsFollowing(followerID, followedID uint) bool {
	key := fmt.Sprintf("%d-%d", followerID, followedID)
	return m.Follows[key]
}

func (m *MockUserRepository) GetFollowedUserIDs(followerID uint) ([]uint, error) {
	if m.ErrGetFollowedUserIDs != nil {
		return nil, m.ErrGetFollowedUserIDs
	}
	var ids []uint
	for key, val := range m.Follows {
		if val {
			var fID, fwID uint
			fmt.Sscanf(key, "%d-%d", &fID, &fwID)
			if fID == followerID {
				ids = append(ids, fwID)
			}
		}
	}
	return ids, nil
}
