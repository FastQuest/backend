package auth

import (
	"errors"
	"time"

	"flashquest/internal/auth/dto"
)

// Domain errors of the auth context. Infrastructure errors (gorm, pgconn) are
// translated into these by the repository, so service and handler never
// depend on the persistence technology.
var (
	ErrNotFound             = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrDuplicatedEmail      = errors.New("duplicated email")
	ErrRoleDomainNotAllowed = errors.New("role domain not allowed")
)

// Credentials is the main entity of the auth domain: the identity data needed
// to authenticate someone and to issue tokens for them.
type Credentials struct {
	UserID       uint
	Name         string
	Email        string
	PasswordHash string
}

// Repository describes the persistence contract of the auth domain.
type Repository interface {
	FindCredentialsByEmail(email string) (*Credentials, error)
	CreateUserWithRole(credentials *Credentials, roleName string) error
	SaveRefreshToken(userID uint, tokenHash string, expiresAt time.Time) error
	GetUserRole(userID uint) (string, error)
}

// Service describes the business contract of the auth domain.
type Service interface {
	Register(req dto.RegisterRequest) (dto.AuthResponse, error)
	Login(req dto.LoginRequest) (dto.AuthResponse, error)
}
