package user

// User is the GORM persistence model for the users table.
// It carries no json tags on purpose: serialization is the job of the DTOs
// in internal/user/dto.
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
