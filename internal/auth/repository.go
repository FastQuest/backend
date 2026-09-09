package auth

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	FindUserByEmail(email string) (*User, error)
	CreateUserWithRole(user *User, roleName string) error
	SaveRefreshToken(userID uint, tokenHash string, expiresAt time.Time) error
	GetUserRole(userID uint) (string, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepositoryWithDB(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) CreateUserWithRole(user *User, roleName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			if isDuplicateKeyError(err) {
				return ErrDuplicatedEmail
			}
			return err
		}

		var role Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			return err
		}

		userRole := UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}
		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormRepository) SaveRefreshToken(userID uint, tokenHash string, expiresAt time.Time) error {
	refreshToken := RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(&refreshToken).Error
}

func (r *GormRepository) GetUserRole(userID uint) (string, error) {
	type roleResult struct {
		Name string
	}

	var role roleResult
	err := r.db.
		Table("roles").
		Select("roles.name").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	return role.Name, nil
}

func isDuplicateKeyError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
