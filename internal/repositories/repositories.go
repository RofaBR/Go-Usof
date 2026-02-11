package repositories

import (
	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/storage/postgres"
	"github.com/RofaBR/Go-Usof/internal/storage/redis"
)

type Repository struct {
	User     domain.UserRepository
	Token    domain.TokenRepository
	Category domain.CategoryRepository
}

func NewRepository(db *postgres.Postgres, rdb *redis.Redis) *Repository {
	return &Repository{
		User:     NewUserRepository(db.Pool),
		Token:    NewTokenRepository(rdb.Client),
		Category: NewCategoryRepository(db.Pool),
	}
}
