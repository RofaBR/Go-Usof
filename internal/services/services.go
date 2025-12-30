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
	mailSenderSvc := NewSMPTSender(config.Sender)
	userSvc := NewUserService(repos.User, tokenSvc, mailSenderSvc)

	return &Service{
		User: userSvc,
	}
}
