package auth

import "time"

// GORM persistence models owned by the auth domain. None of them carry json
// tags: serialization belongs to the DTOs in internal/auth/dto.

// User is the auth view of the users table, limited to the columns this
// domain needs to authenticate someone.
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

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
