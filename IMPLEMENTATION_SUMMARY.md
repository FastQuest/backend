# FastQuest Auth RS256 Implementation - Complete

## Summary
Implemented secure JWT authentication system with RS256 signing, refresh token persistence, and automatic role assignment based on email domain.

## Commits (7 tasks)

### Task 1: Standardized API Responses ✅
- **Commit:** `50bc102` → `a338aa1` (fixed encoding safety)
- `pkg/apiresp/response.go` - standardized error/JSON helpers
- Enforces consistent response format across all endpoints

### Task 2: Security Primitives ✅
- **Commit:** `18bdd72` → `71ecd99` (fixed JWT claim contract)
- `pkg/security/password/` - bcrypt hashing
- `pkg/security/token/` - random token generation with hash verification
- `pkg/security/jwt/` - RS256 signing/parsing with required claim validation

### Task 3: Database Schema & Models ✅
- **Commit:** `3ed4924`
- Roles table with seed data (Aluno, Professor, Admin)
- User-roles junction table for flexible role assignment
- Refresh tokens table with hash-only storage and TTL

### Task 4a: Auth Service (Register/Login) ✅
- **Commit:** `55a7f02`
- RegisterRequest/LoginRequest/AuthResponse DTOs
- Service layer with role mapping by email domain
- Repository abstraction over GORM

### Task 4b: Input Validation Hardening ✅
- **Commit:** `a5eabea`
- Email format validation (RFC5322 via net/mail.ParseAddress)
- Password minimum length enforcement (6+ chars)
- User enumeration protection (same error for unknown user and wrong password)
- Comprehensive validation tests for all input scenarios

### Task 5: HTTP Handlers + Routes ✅
- **Commit:** `0618399`
- RegisterHandler and LoginHandler HTTP functions
- Routes: `POST /api/auth/register` and `POST /api/auth/login`
- Standardized error responses with correct HTTP status codes

### Task 6: JWT Middleware ✅
- **Commit:** `d7a227c`
- RequireAuth middleware as http.Handler wrapper
- Bearer token extraction and RS256 signature validation
- Context injection (user_id, role) for downstream handlers

## API Contracts

### Register Endpoint
```
POST /api/auth/register
Content-Type: application/json

Request:
{
  "name": "John Doe",
  "email": "john@sempreceub.com",
  "password": "securepass123"
}

Response (200):
{
  "access_token": "eyJhbGc...",
  "expires_in": 259200,
  "user_id": 123
}

Error Responses:
- 400 VALIDATION_ERROR: Invalid email format, empty password, etc.
- 409 EMAIL_ALREADY_EXISTS: User already registered
- 422 ROLE_DOMAIN_NOT_ALLOWED: Email domain not supported
```

### Login Endpoint
```
POST /api/auth/login
Content-Type: application/json

Request:
{
  "email": "john@sempreceub.com",
  "password": "securepass123"
}

Response (200):
{
  "access_token": "eyJhbGc...",
  "expires_in": 259200,
  "user_id": 123
}

Error Responses:
- 400 VALIDATION_ERROR: Invalid email format, missing password
- 401 INVALID_CREDENTIALS: Unknown email or wrong password
```

### Protected Endpoints (via Middleware)
```
GET /api/protected-resource
Authorization: Bearer <access_token>

Middleware validates token and injects:
- Context["user_id"]: uint (JWT sub claim)
- Context["role"]: string (JWT role claim)

Error Responses:
- 401 MISSING_TOKEN: Authorization header absent
- 401 INVALID_TOKEN_FORMAT: Bearer format incorrect
- 401 INVALID_TOKEN: Signature/expiry validation failed
```

## Email Domain to Role Mapping
- `@sempreceub.com` → Aluno
- `@ceub.edu.br` → Professor
- Other domains → 422 ROLE_DOMAIN_NOT_ALLOWED
- Admin role assignment requires manual DB entry (future scope)

## Security Features
- Passwords hashed with bcrypt (cost 10)
- Refresh tokens: opaque random + hashed in DB (prevents token exfiltration)
- JWT: RS256 with 72-hour expiration
- Claims: `sub` (user ID), `exp` (expiry), `role` (group)
- Input validation: email format, password strength, name required
- User enumeration protection: same 401 error for unknown user and wrong password
- SQL injection: GORM parameterized queries throughout
- XSS mitigation: structured JSON responses only

## Test Coverage
- 12 integration tests for auth service
- Validation edge cases (invalid emails, short passwords, empty inputs)
- Role mapping by domain (both valid and unsupported domains)
- Token issuance and persistence
- User enumeration protection

## Build & Test Status
```bash
✅ go test ./...     # All tests passing
✅ go build ./...    # All packages building
✅ Worktree safety   # .worktrees/ in .gitignore
```

## Next Steps (Out of Scope for This Phase)
1. POST /api/auth/refresh - token refresh with refresh token rotation
2. POST /api/auth/logout - refresh token revocation
3. Role-based access control middleware (enforce roles on protected endpoints)
4. Role management endpoints (assign/revoke roles)
5. Integration tests with actual PostgreSQL
6. OpenAPI/Swagger documentation

## Branch: feat/auth-rs256
All work isolated in worktree `.worktrees/auth-rs256` on branch `feat/auth-rs256`
Ready for review and integration to main.
