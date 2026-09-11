package auth

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"flashquest/internal/auth/dto"
	jwtsec "flashquest/pkg/security/jwt"
	"flashquest/pkg/security/password"
	tokensec "flashquest/pkg/security/token"
)

const (
	accessTokenTTL       = 72 * time.Hour
	accessTokenExpiresIn = int64(259200)
)

type service struct {
	repository    Repository
	privateKeyPEM string
}

// NewService builds the default implementation of Service.
func NewService(repository Repository, privateKeyPEM string) Service {
	return &service{
		repository:    repository,
		privateKeyPEM: privateKeyPEM,
	}
}

func (s *service) Register(req dto.RegisterRequest) (dto.AuthResponse, error) {
	// Validate inputs
	if err := validateRegisterRequest(req); err != nil {
		return dto.AuthResponse{}, err
	}

	normalizedEmail := normalizeEmail(req.Email)
	roleName, err := resolveRole(normalizedEmail)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	existingCredentials, err := s.repository.FindCredentialsByEmail(normalizedEmail)
	if err == nil && existingCredentials != nil {
		return dto.AuthResponse{}, ErrDuplicatedEmail
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return dto.AuthResponse{}, err
	}

	passwordHash, err := password.Hash(req.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	credentials := &Credentials{
		Name:         req.Name,
		Email:        normalizedEmail,
		PasswordHash: passwordHash,
	}
	if err := s.repository.CreateUserWithRole(credentials, roleName); err != nil {
		if errors.Is(err, ErrDuplicatedEmail) {
			return dto.AuthResponse{}, ErrDuplicatedEmail
		}
		return dto.AuthResponse{}, err
	}

	return s.issueTokens(credentials.UserID, roleName)
}

func (s *service) Login(req dto.LoginRequest) (dto.AuthResponse, error) {
	// Validate inputs
	if err := validateLoginRequest(req); err != nil {
		return dto.AuthResponse{}, err
	}

	normalizedEmail := normalizeEmail(req.Email)
	credentials, err := s.repository.FindCredentialsByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return dto.AuthResponse{}, ErrInvalidCredentials
		}
		return dto.AuthResponse{}, err
	}

	if err := password.Compare(credentials.PasswordHash, req.Password); err != nil {
		return dto.AuthResponse{}, ErrInvalidCredentials
	}

	roleName, err := s.repository.GetUserRole(credentials.UserID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return dto.AuthResponse{}, ErrInvalidCredentials
		}
		return dto.AuthResponse{}, err
	}

	return s.issueTokens(credentials.UserID, roleName)
}

func (s *service) issueTokens(userID uint, roleName string) (dto.AuthResponse, error) {
	accessToken, err := jwtsec.SignAccessToken(jwtsec.AuthClaims{
		UserID: userID,
		Role:   roleName,
	}, accessTokenTTL, s.privateKeyPEM)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	refreshToken, refreshHash, err := tokensec.GenerateAndHash()
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if err := s.repository.SaveRefreshToken(userID, refreshHash, time.Now().Add(accessTokenTTL)); err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTokenExpiresIn,
		UserID:       userID,
	}, nil
}

func resolveRole(email string) (string, error) {
	switch {
	case strings.HasSuffix(email, "@sempreceub.com"):
		return "Aluno", nil
	case strings.HasSuffix(email, "@ceub.edu.br"):
		return "Professor", nil
	default:
		return "", ErrRoleDomainNotAllowed
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateRegisterRequest(req dto.RegisterRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if req.Password == "" || len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func validateLoginRequest(req dto.LoginRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func validateEmail(email string) error {
	normalized := normalizeEmail(email)
	if normalized == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(normalized)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}
