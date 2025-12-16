package router

import (
	"github.com/RofaBR/Go-Usof/internal/handler"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(log *logger.Logger) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	healthHandler := handler.NewHealthHandler(log)

	if gin.Mode() != gin.ReleaseMode {
		registerHealthRoutes(router, healthHandler)
	}

	registerAuthRoutes(api)

	return router
}

func registerHealthRoutes(router *gin.Engine, h *handler.HealthHandler) {
	router.GET("/ping", h.Ping)
}

func registerAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register")
		auth.POST("/login")
		auth.POST("/logout")
	}
}
