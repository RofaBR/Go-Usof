package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RofaBR/Go-Usof/internal/domain"
	goredis "github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *goredis.Client
}

func NewTokenRepository(client *goredis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (t *TokenRepository) StoreRefreshToken(ctx context.Context, metadata *domain.RefreshTokenMetadata, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s:%s", metadata.UserID, metadata.JTI)

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return t.client.Set(ctx, key, data, ttl).Err()
}

func (t *TokenRepository) GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshTokenMetadata, error) {
	pattern := fmt.Sprintf("refresh:*:%s", jti)
	keys, err := t.client.Keys(ctx, pattern).Result()

	if len(keys) == 0 {
		return nil, nil
	}

	data, err := t.client.Get(ctx, keys[0]).Result()
	if err == goredis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	var metadata domain.RefreshTokenMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

func (t *TokenRepository) DeleteRefreshToken(ctx context.Context, jti string) error {
	pattern := fmt.Sprintf("refresh:*:%s", jti)
	keys, err := t.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get refresh token keys: %w", err)
	}
	if len(keys) > 0 {
		return t.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (t *TokenRepository) ExtendRefreshTokenTTL(ctx context.Context, jti string, ttl time.Duration) error {
	pattern := fmt.Sprintf("refresh:*:%s", jti)
	keys, err := t.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find token keys: %w", err)
	}
	if len(keys) == 0 {
		return fmt.Errorf("refresh token not found")
	}
	return t.client.Expire(ctx, keys[0], ttl).Err()
}

func (t *TokenRepository) StoreVerificationToken(ctx context.Context, metadata *domain.VerificationTokenMetadata, ttl time.Duration) error {
	key := fmt.Sprintf("verification:%s", metadata.Token)

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal verification token metadata: %w", err)
	}

	return t.client.Set(ctx, key, data, ttl).Err()
}

func (t *TokenRepository) GetVerificationToken(ctx context.Context, token string) (*domain.VerificationTokenMetadata, error) {
	key := fmt.Sprintf("verification:%s", token)

	data, err := t.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get verification token: %w", err)
	}

	var metadata domain.VerificationTokenMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal verification token metadata: %w", err)
	}

	return &metadata, nil
}

func (t *TokenRepository) DeleteVerificationToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("verification:%s", token)
	return t.client.Del(ctx, key).Err()
}
