package middleware

import (
	"net/http"
	"strings"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/gin-gonic/gin"
)

const ClaimsKey = "claims"

func AuthMiddleware(tokenService *services.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format: expected 'Bearer <token>'"})
			c.Abort()
			return
		}
		if parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format: expected 'Bearer'"})
			c.Abort()
			return
		}
		claims, err := tokenService.ValidateAccessToken(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get(ClaimsKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := raw.(*domain.TokenClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims type"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
