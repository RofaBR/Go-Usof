package services

import (
	"github.com/RofaBR/Go-Usof/internal/config"
	"github.com/RofaBR/Go-Usof/internal/repositories"
	"github.com/RofaBR/Go-Usof/pkg/logger"
)

type Service struct {
	User  *UserService
	Token *TokenService
	Email *SMTPSender
	Image *CloudinaryService
}

func NewServices(log *logger.Logger, repos *repositories.Repository, config *config.Config) *Service {
	tokenSvc := NewTokenService(repos.Token, config.JWT)
	emailSvc := NewSMTPSender(config.Sender)
	cloudinarySvc := NewCloudinaryService(config.CloudinaryURL)
	userSvc := NewUserService(repos.User, log)

	return &Service{
		User:  userSvc,
		Token: tokenSvc,
		Email: emailSvc,
		Image: cloudinarySvc,
	}
}
