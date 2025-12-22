package handler

import (
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/pkg/logger"
)

type Handler struct {
	Health *HealthHandler
	Auth   *AuthHandler
}

func NewHandler(log *logger.Logger, svc *services.Service) *Handler {
	return &Handler{
		Health: NewHealthHandler(log),
		Auth:   NewAuthHandler(svc.User),
	}
}
