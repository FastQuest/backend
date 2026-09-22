package user

import (
	"errors"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

// NewRepository builds the GORM backed implementation of Repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetUserByID(userID uint) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
