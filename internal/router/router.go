package router

import (
	"github.com/RofaBR/Go-Usof/internal/handler"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(log *logger.Logger) *gin.Engine {
	router := gin.Default()

	healthHandler := handler.NewHealthHandler(log)

	router.GET("/ping", healthHandler.Ping)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register")
			auth.POST("/login")
			auth.POST("/logout")
		}
	}

	return router
}
