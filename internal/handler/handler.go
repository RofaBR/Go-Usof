package handler

import (
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/pkg/logger"
)

type Handler struct {
	Health *HealthHandler
	Auth   *AuthHandler
	User   *UserHandler
}

func NewHandler(log *logger.Logger, svc *services.Service) *Handler {
	return &Handler{
		Health: NewHealthHandler(log),
		Auth:   NewAuthHandler(svc.User, svc.Token, svc.Email, log),
		User:   NewUserHandler(svc.User, svc.Image, svc.Token, log),
	}
}
