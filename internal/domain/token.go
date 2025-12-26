package domain

import (
	"context"
	"time"
)

type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`         // Access token TTL in seconds
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // Refresh token TTL in seconds
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	JTI    string `json:"jti"`
	Type   string `json:"type"`
}

type RefreshTokenMetadata struct {
	UserID           string    `json:"user_id"`
	JTI              string    `json:"jti"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
	AbsoluteExpireAt time.Time `json:"absolute_expire_at"`
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, metadata *RefreshTokenMetadata, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, jti string) (*RefreshTokenMetadata, error)
	//ExtendRefreshTokenTTL(ctx context.Context, jti string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, jti string) error
}

type TokenService interface {
	GenerateTokenPair(ctx context.Context, user *User) (*TokenPair, error)

	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*TokenClaims, error)

	//RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenPair, error)

	//RevokeToken(ctx context.Context, accessToken, refreshToken string) error
}
