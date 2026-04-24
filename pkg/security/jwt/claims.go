package jwt

type AuthClaims struct {
	UserID    uint   `json:"sub"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat,omitempty"`
}
