package repository

import (
	"time"

	"realworld-endpoints/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error

	// Follow System
	Follow(followerID, followedID uint) error
	Unfollow(followerID, followedID uint) error
	IsFollowing(followerID, followedID uint) bool
	GetFollowedUserIDs(followerID uint) ([]uint, error)
	BatchIsFollowing(followerID uint, followedIDs []uint) (map[uint]bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Follow(followerID, followedID uint) error {
	follow := models.UserFollow{
		FollowerID: followerID,
		FollowedID: followedID,
		CreatedAt:  time.Now(),
	}
	return r.db.Where(models.UserFollow{FollowerID: followerID, FollowedID: followedID}).
		FirstOrCreate(&follow).Error
}

func (r *userRepository) Unfollow(followerID, followedID uint) error {
	return r.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Delete(&models.UserFollow{}).Error
}

func (r *userRepository) IsFollowing(followerID, followedID uint) bool {
	if followerID == 0 || followedID == 0 {
		return false
	}
	var count int64
	r.db.Model(&models.UserFollow{}).
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Count(&count)
	return count > 0
}

func (r *userRepository) GetFollowedUserIDs(followerID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.UserFollow{}).
		Where("follower_id = ?", followerID).
		Pluck("followed_id", &ids).Error
	return ids, err
}

func (r *userRepository) BatchIsFollowing(followerID uint, followedIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if followerID == 0 || len(followedIDs) == 0 {
		return result, nil
	}
	var followingIDs []uint
	err := r.db.Model(&models.UserFollow{}).
		Where("follower_id = ? AND followed_id IN ?", followerID, followedIDs).
		Pluck("followed_id", &followingIDs).Error
	if err != nil {
		return nil, err
	}
	for _, id := range followingIDs {
		result[id] = true
	}
	return result, nil
}

