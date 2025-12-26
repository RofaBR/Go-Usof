package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RofaBR/Go-Usof/internal/config"
	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type TokenService struct {
	repo   domain.TokenRepository
	config config.JWTConfig
}

func NewTokenService(repo domain.TokenRepository, cfg config.JWTConfig) *TokenService {
	return &TokenService{
		repo:   repo,
		config: cfg,
	}
}

func (t *TokenService) GenerateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	accessJTI := uuid.New().String()
	accessToken, err := t.generateAccessToken(user, accessJTI)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshStr := uuid.New().String()
	refreshToken, err := t.generateRefreshToken(user, refreshStr)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	metadata := domain.RefreshTokenMetadata{
		UserID:           strconv.Itoa(user.ID),
		JTI:              refreshStr,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(time.Duration(t.config.RefreshTTL) * 24 * time.Hour),
		AbsoluteExpireAt: time.Now().Add(30 * 24 * time.Hour),
	}

	ttl := time.Duration(t.config.RefreshTTL) * 24 * time.Hour
	if err := t.repo.StoreRefreshToken(ctx, &metadata, ttl); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        int64(t.config.AccessTTL * 60),
		RefreshExpiresIn: int64(t.config.RefreshTTL * 24 * 60 * 60),
	}, nil
}

func (t *TokenService) ValidateAccessToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	claims, err := t.parseToken(token, t.config.AccessSecret)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}
func (t *TokenService) ValidateRefreshToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	claims, err := t.parseToken(token, t.config.RefreshSecret)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	metadata, err := t.repo.GetRefreshToken(ctx, claims.JTI)
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %w", err)
	}
	if metadata == nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	}
	return claims, nil
}

func (t *TokenService) parseToken(tokenString, secret string) (*domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID := ""
	if val, ok := claims["user_id"].(float64); ok {
		userID = strconv.Itoa(int(val))
	}

	return &domain.TokenClaims{
		UserID: userID,
		Email:  getStringClaim(claims, "email"),
		Role:   getStringClaim(claims, "role"),
		JTI:    getStringClaim(claims, "jti"),
		Type:   getStringClaim(claims, "type"),
	}, nil
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}

func (t *TokenService) generateAccessToken(user *domain.User, jti string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"type":    "access",
		"exp":     time.Now().Add(time.Duration(t.config.AccessTTL) * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
		"jti":     jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.config.AccessSecret))
}

func (t *TokenService) generateRefreshToken(user *domain.User, jti string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Duration(t.config.RefreshTTL) * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"jti":     jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.config.RefreshSecret))
}

func (t *TokenService) RevokeToken(ctx context.Context, refreshToken string) error {
	claims, err := t.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}
	return t.repo.DeleteRefreshToken(ctx, claims.JTI)
}
