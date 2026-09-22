package auth

import (
	"errors"
	"time"

	"flashquest/internal/user"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

// NewRepositoryWithDB builds the GORM backed implementation of Repository.
func NewRepositoryWithDB(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) FindCredentialsByEmail(email string) (*Credentials, error) {
	var user user.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toCredentials(user), nil
}

func (r *gormRepository) CreateUserWithRole(credentials *Credentials, roleName string) error {
	user := user.User{
		Name:         credentials.Name,
		Email:        credentials.Email,
		PasswordHash: credentials.PasswordHash,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			if isDuplicateKeyError(err) {
				return ErrDuplicatedEmail
			}
			return err
		}

		var role Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
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
	if err != nil {
		return err
	}

	credentials.UserID = user.ID
	return nil
}

func (r *gormRepository) SaveRefreshToken(userID uint, tokenHash string, expiresAt time.Time) error {
	refreshToken := RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(&refreshToken).Error
}

func (r *gormRepository) GetUserRole(userID uint) (string, error) {
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
			return "", ErrNotFound
		}
		return "", err
	}

	return role.Name, nil
}

func toCredentials(user user.User) *Credentials {
	return &Credentials{
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
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
