package router

import (
	"github.com/RofaBR/Go-Usof/internal/handler"
	"github.com/RofaBR/Go-Usof/internal/middleware"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(log *logger.Logger, h *handler.Handler, authMW gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	if gin.Mode() != gin.ReleaseMode {
		registerHealthRoutes(router, h)
	}

	registerAuthRoutes(api, h)
	registerUserRoutes(api, h)
	registerCategoryRoutes(api, h, authMW)

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
}

func registerUserRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	user := rg.Group("/user")
	{
		user.POST("/register", h.Auth.Register)
	}
}

func registerCategoryRoutes(rg *gin.RouterGroup, h *handler.Handler, authMW gin.HandlerFunc) {
	category := rg.Group("/category")
	category.Use(authMW)
	{
		category.POST("/create", middleware.RoleMiddleware("admin"), h.Category.Create)
	}
}
