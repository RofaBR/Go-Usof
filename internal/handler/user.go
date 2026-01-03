package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  *services.UserService
	imageService *services.CloudinaryService
	tokenService *services.TokenService
	log          *logger.Logger
}

func NewUserHandler(userService *services.UserService, imageService *services.CloudinaryService, tokenService *services.TokenService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService:  userService,
		imageService: imageService,
		tokenService: tokenService,
		log:          log,
	}
}

func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling update avatar request")

	token, err := extractToken(c)
	if err != nil {
		h.log.Warn("failed to extract token from request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	params, err := h.tokenService.ValidateAccessToken(c.Request.Context(), token)
	if err != nil {
		h.log.Warn("failed to validate access token", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userIDStr := params.UserID

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.log.Error("invalid user ID format", "userId", userIDStr, "error", err)
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		h.log.Warn("failed to parse avatar file", "error", err)
		c.JSON(400, gin.H{"error": "Avatar file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Error("failed to open avatar file", "error", err)
		c.JSON(400, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	avatarURL, err := h.imageService.UploadAvatar(ctx, file, userIDStr)
	if err != nil {
		h.log.Error("failed to upload avatar", "userId", userID, "error", err)
		c.JSON(500, gin.H{"error": "Failed to upload image"})
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		h.log.Error("failed to get user", "userId", userID, "error", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if user == nil {
		h.log.Warn("user not found", "userId", userID)
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user.Avatar = avatarURL

	err = h.userService.Update(ctx, user)
	if err != nil {
		h.log.Error("failed to update user avatar", "userId", userID, "error", err)
		c.JSON(500, gin.H{"error": "Failed to update avatar"})
		return
	}

	h.log.Info("avatar updated successfully", "userId", userID, "avatarURL", avatarURL)
	c.JSON(200, gin.H{
		"message":    "Avatar updated successfully",
		"avatar_url": avatarURL,
	})
}

func extractToken(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid token format: expected 'Bearer <token>'")
	}

	if parts[0] != "Bearer" {
		return "", errors.New("invalid token format: expected Bearer scheme")
	}

	return parts[1], nil
}
