package models

import "time"

type RefreshToken struct {
ID        uint       `gorm:"primaryKey"`
UserID    uint       `gorm:"not null"`
TokenHash string     `gorm:"not null;unique"`
ExpiresAt time.Time  `gorm:"not null"`
RevokedAt *time.Time
CreatedAt time.Time
}

func (RefreshToken) TableName() string {
return "refresh_tokens"
}
