package repositories

import (
	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/storage/postgres"
)

type Repository struct {
	User domain.UserRepository
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		User: NewUserRepository(db.Pool),
	}
}
