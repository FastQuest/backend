package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
	"time"

	"flashquest/pkg/models"
	jwtsec "flashquest/pkg/security/jwt"
	"flashquest/pkg/security/password"
	tokensec "flashquest/pkg/security/token"
)

func TestResolveRoleByEmailDomain(t *testing.T) {
	role, err := resolveRole("x@sempreceub.com")
	if err != nil {
		t.Fatalf("expected no error for aluno domain, got %v", err)
	}
	if role != "Aluno" {
		t.Fatalf("expected Aluno role, got %q", role)
	}

	role, err = resolveRole("x@ceub.edu.br")
	if err != nil {
		t.Fatalf("expected no error for professor domain, got %v", err)
	}
	if role != "Professor" {
		t.Fatalf("expected Professor role, got %q", role)
	}

	if _, err := resolveRole("x@gmail.com"); !errors.Is(err, ErrRoleDomainNotAllowed) {
		t.Fatalf("expected ErrRoleDomainNotAllowed, got %v", err)
	}
}

func TestRegisterCreatesUserRoleAndTokens(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)
	repo := &stubRepository{
		findUserByEmailFn: func(_ string) (*models.User, error) {
			return nil, ErrUserNotFound
		},
		createUserWithRoleFn: func(user *models.User, roleName string) error {
			user.ID = 123
			if roleName != "Aluno" {
				t.Fatalf("expected role Aluno, got %q", roleName)
			}
			if user.Email != "user@sempreceub.com" {
				t.Fatalf("expected normalized email user@sempreceub.com, got %q", user.Email)
			}
			if user.PasswordHash == "password123" {
				t.Fatal("expected hashed password")
			}
			if err := password.Compare(user.PasswordHash, "password123"); err != nil {
				t.Fatalf("expected bcrypt hash matching raw password, got %v", err)
			}
			return nil
		},
		getUserRoleFn: func(uint) (string, error) {
			t.Fatal("did not expect GetUserRole on register")
			return "", nil
		},
	}

	service := NewService(repo, privateKeyPEM)

	got, err := service.Register(RegisterRequest{
		Name:     "User",
		Email:    "  USER@SEMPRECEUB.COM  ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.UserID != 123 {
		t.Fatalf("expected user id 123, got %d", got.UserID)
	}
	if got.ExpiresIn != 259200 {
		t.Fatalf("expected expires_in 259200, got %d", got.ExpiresIn)
	}
	if strings.TrimSpace(got.AccessToken) == "" {
		t.Fatal("expected non-empty access token")
	}
	if strings.TrimSpace(got.RefreshToken) == "" {
		t.Fatal("expected non-empty refresh token")
	}

	if repo.savedRefreshUserID != 123 {
		t.Fatalf("expected refresh token persisted for user 123, got %d", repo.savedRefreshUserID)
	}
	if repo.savedRefreshTokenHash == "" {
		t.Fatal("expected refresh token hash persisted")
	}
	if !tokensec.VerifyHash(got.RefreshToken, repo.savedRefreshTokenHash) {
		t.Fatal("expected persisted refresh hash to match returned refresh token")
	}
	if repo.savedRefreshExpiresAt.Before(time.Now().Add(71 * time.Hour)) || repo.savedRefreshExpiresAt.After(time.Now().Add(73*time.Hour)) {
		t.Fatalf("expected refresh expiration around 72h, got %s", repo.savedRefreshExpiresAt.Sub(time.Now()))
	}

	claims, err := jwtsec.ParseAccessToken(got.AccessToken, publicKeyPEM)
	if err != nil {
		t.Fatalf("expected valid access token, got %v", err)
	}
	if claims.UserID != 123 {
		t.Fatalf("expected token sub 123, got %d", claims.UserID)
	}
	if claims.Role != "Aluno" {
		t.Fatalf("expected token role Aluno, got %q", claims.Role)
	}
}

func TestRegisterReturnsDuplicatedEmail(t *testing.T) {
	service := NewService(&stubRepository{
		findUserByEmailFn: func(_ string) (*models.User, error) {
			return &models.User{ID: 1}, nil
		},
	}, "key")

	_, err := service.Register(RegisterRequest{
		Name:     "User",
		Email:    "user@sempreceub.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrDuplicatedEmail) {
		t.Fatalf("expected ErrDuplicatedEmail, got %v", err)
	}
}

func TestRegisterReturnsRoleDomainNotAllowed(t *testing.T) {
	service := NewService(&stubRepository{}, "key")

	_, err := service.Register(RegisterRequest{
		Name:     "User",
		Email:    "user@gmail.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrRoleDomainNotAllowed) {
		t.Fatalf("expected ErrRoleDomainNotAllowed, got %v", err)
	}
}

func TestLoginIssuesTokensAndPersistsRefreshHash(t *testing.T) {
	privateKeyPEM, publicKeyPEM := generateTestKeyPair(t)
	passwordHash, err := password.Hash("password123")
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
	}

	repo := &stubRepository{
		findUserByEmailFn: func(email string) (*models.User, error) {
			if email != "prof@ceub.edu.br" {
				t.Fatalf("expected normalized email prof@ceub.edu.br, got %q", email)
			}
			return &models.User{ID: 321, Email: email, PasswordHash: passwordHash}, nil
		},
		getUserRoleFn: func(userID uint) (string, error) {
			if userID != 321 {
				t.Fatalf("expected user id 321 for role lookup, got %d", userID)
			}
			return "Professor", nil
		},
	}
	service := NewService(repo, privateKeyPEM)

	got, err := service.Login(LoginRequest{
		Email:    "  PROF@CEUB.EDU.BR ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.UserID != 321 {
		t.Fatalf("expected user id 321, got %d", got.UserID)
	}
	if got.ExpiresIn != 259200 {
		t.Fatalf("expected expires_in 259200, got %d", got.ExpiresIn)
	}
	if strings.TrimSpace(got.AccessToken) == "" || strings.TrimSpace(got.RefreshToken) == "" {
		t.Fatal("expected non-empty access and refresh token")
	}
	if repo.savedRefreshUserID != 321 {
		t.Fatalf("expected refresh token persisted for user 321, got %d", repo.savedRefreshUserID)
	}
	if !tokensec.VerifyHash(got.RefreshToken, repo.savedRefreshTokenHash) {
		t.Fatal("expected persisted refresh hash to match returned refresh token")
	}

	claims, err := jwtsec.ParseAccessToken(got.AccessToken, publicKeyPEM)
	if err != nil {
		t.Fatalf("expected valid access token, got %v", err)
	}
	if claims.UserID != 321 || claims.Role != "Professor" {
		t.Fatalf("expected claims (321, Professor), got (%d, %s)", claims.UserID, claims.Role)
	}
}

func TestLoginReturnsInvalidCredentials(t *testing.T) {
	service := NewService(&stubRepository{
		findUserByEmailFn: func(_ string) (*models.User, error) {
			return nil, ErrUserNotFound
		},
	}, "key")

	_, err := service.Login(LoginRequest{
		Email:    "user@sempreceub.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

type stubRepository struct {
	findUserByEmailFn  func(email string) (*models.User, error)
	createUserWithRoleFn func(user *models.User, roleName string) error
	saveRefreshTokenFn func(userID uint, tokenHash string, expiresAt time.Time) error
	getUserRoleFn      func(userID uint) (string, error)

	savedRefreshUserID   uint
	savedRefreshTokenHash string
	savedRefreshExpiresAt time.Time
}

func (s *stubRepository) FindUserByEmail(email string) (*models.User, error) {
	if s.findUserByEmailFn == nil {
		return nil, ErrUserNotFound
	}
	return s.findUserByEmailFn(email)
}

func (s *stubRepository) CreateUserWithRole(user *models.User, roleName string) error {
	if s.createUserWithRoleFn == nil {
		return nil
	}
	return s.createUserWithRoleFn(user, roleName)
}

func (s *stubRepository) SaveRefreshToken(userID uint, tokenHash string, expiresAt time.Time) error {
	s.savedRefreshUserID = userID
	s.savedRefreshTokenHash = tokenHash
	s.savedRefreshExpiresAt = expiresAt
	if s.saveRefreshTokenFn == nil {
		return nil
	}
	return s.saveRefreshTokenFn(userID, tokenHash, expiresAt)
}

func (s *stubRepository) GetUserRole(userID uint) (string, error) {
	if s.getUserRoleFn == nil {
		return "", nil
	}
	return s.getUserRoleFn(userID)
}

func generateTestKeyPair(t *testing.T) (string, string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateDER})

	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to encode RSA public key: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})

	return string(privatePEM), string(publicPEM)
}
<<<<<<< HEAD
=======

func TestRegisterValidateEmail(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	repo := &stubRepository{}
	svc := NewService(repo, privateKeyPEM)

	// Test empty email
	_, err := svc.Register(RegisterRequest{
		Name:     "User",
		Email:    "",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for empty email")
	}

	// Test invalid email format
	_, err = svc.Register(RegisterRequest{
		Name:     "User",
		Email:    "notanemail",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for invalid email format")
	}

	// Test email without local part
	_, err = svc.Register(RegisterRequest{
		Name:     "User",
		Email:    "@ceub.edu.br",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for email without local part")
	}
}

func TestRegisterValidatePassword(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	repo := &stubRepository{}
	svc := NewService(repo, privateKeyPEM)

	// Test empty password
	_, err := svc.Register(RegisterRequest{
		Name:     "User",
		Email:    "user@sempreceub.com",
		Password: "",
	})
	if err == nil {
		t.Fatal("expected error for empty password")
	}

	// Test password too short
	_, err = svc.Register(RegisterRequest{
		Name:     "User",
		Email:    "user@sempreceub.com",
		Password: "short",
	})
	if err == nil {
		t.Fatal("expected error for password < 6 chars")
	}
}

func TestRegisterValidateName(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	repo := &stubRepository{}
	svc := NewService(repo, privateKeyPEM)

	// Test empty name
	_, err := svc.Register(RegisterRequest{
		Name:     "",
		Email:    "user@sempreceub.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestLoginValidateEmail(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	repo := &stubRepository{}
	svc := NewService(repo, privateKeyPEM)

	// Test empty email
	_, err := svc.Login(LoginRequest{
		Email:    "",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for empty email")
	}

	// Test invalid email format
	_, err = svc.Login(LoginRequest{
		Email:    "notanemail",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for invalid email format")
	}
}

func TestLoginValidatePassword(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	repo := &stubRepository{}
	svc := NewService(repo, privateKeyPEM)

	// Test empty password
	_, err := svc.Login(LoginRequest{
		Email:    "user@sempreceub.com",
		Password: "",
	})
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestLoginUserEnumeration(t *testing.T) {
	privateKeyPEM, _ := generateTestKeyPair(t)
	repo := &stubRepository{
		findUserByEmailFn: func(_ string) (*models.User, error) {
			return nil, ErrUserNotFound
		},
	}
	svc := NewService(repo, privateKeyPEM)

	// Both non-existent and wrong password should return same generic error
	_, err := svc.Login(LoginRequest{
		Email:    "unknown@sempreceub.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
>>>>>>> feat/auth-rs256
