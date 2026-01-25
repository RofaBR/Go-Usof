package router

import (
	"github.com/RofaBR/Go-Usof/internal/handler"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(log *logger.Logger, h *handler.Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	if gin.Mode() != gin.ReleaseMode {
		registerHealthRoutes(router, h)
	}

	registerAuthRoutes(api, h)

	return router
}

func registerHealthRoutes(router *gin.Engine, h *handler.Handler) {
	router.GET("/ping", h.Health.Ping)
}

func registerAuthRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/logout", h.Auth.Logout)
		auth.GET("/verify", h.Auth.VerifyEmail)
		auth.GET("/google", h.OAuth2.GoogleLogin)
		auth.GET("/google/callback", h.OAuth2.GoogleCallback)

		auth.POST("/refresh", h.Auth.Refresh)
	}

	user := rg.Group("/user")
	{
		user.POST("/upload/avatar", h.User.UpdateAvatar)
	}
}
