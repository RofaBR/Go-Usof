package services

import (
	"github.com/RofaBR/Go-Usof/internal/repositories"
	"github.com/RofaBR/Go-Usof/pkg/logger"
)

type Service struct {
	User *UserService
}

func NewServices(log *logger.Logger, repos *repositories.Repository) *Service {
	userSvc := NewUserService(repos.User)

	return &Service{
		User: userSvc,
	}
}
