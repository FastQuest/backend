package jwt

type AuthClaims struct {
	UserID    uint   `json:"user_id"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
}
