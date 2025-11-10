package handler

import (
	"net/http"
	"time"

	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	logger *logger.Logger
}

func NewHealthHandler(logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func (h *HealthHandler) Ping(c *gin.Context) {
	h.logger.Debug("ping endpoint called")

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "go-usof",
	}

	c.JSON(http.StatusOK, response)
}
