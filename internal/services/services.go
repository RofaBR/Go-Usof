package services

import (
	"github.com/RofaBR/Go-Usof/internal/config"
	"github.com/RofaBR/Go-Usof/internal/repositories"
	"github.com/RofaBR/Go-Usof/pkg/logger"
)

type Service struct {
	User *UserService
}

func NewServices(log *logger.Logger, repos *repositories.Repository, config *config.Config) *Service {
	tokenSvc := NewTokenService(repos.Token, config.JWT)
	userSvc := NewUserService(repos.User, tokenSvc)

	return &Service{
		User: userSvc,
	}
}
