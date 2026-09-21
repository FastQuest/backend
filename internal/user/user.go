package user

import (
	"errors"

	"flashquest/internal/user/dto"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

// ErrNotFound is the domain level translation of a missing user. The
// repository is responsible for converting infrastructure errors (such as
// gorm.ErrRecordNotFound) into this error so the upper layers never depend on
// the persistence technology.
var ErrNotFound = errors.New("user not found")

// Repository describes the persistence contract of the user domain.
// The main entity of this domain is User (see models.go).
type Repository interface {
	GetUserByID(userID uint) (*User, error)
}

// Service describes the business contract of the user domain.
type Service interface {
	GetCurrentUser(userID uint) (dto.Response, error)
}
