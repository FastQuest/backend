package auth

import "time"

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null;unique"`
}

func (Role) TableName() string {
	return "roles"
}

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	RoleID uint `gorm:"not null"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
