package handler

import (
	"net/http"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/dto/request"
	"github.com/RofaBR/Go-Usof/internal/dto/response"
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService  *services.UserService
	tokenService *services.TokenService
	emailService *services.SMTPSender
	log          *logger.Logger
}

func NewAuthHandler(userService *services.UserService, tokenService *services.TokenService, emailService *services.SMTPSender, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
		emailService: emailService,
		log:          log,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling user registration request")

	var req request.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid registration request", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     "user",
	}

	err := h.userService.Create(ctx, user)
	if err != nil {
		h.log.Error("failed to create user", "email", req.Email, "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.tokenService.GenerateVerificationToken(ctx, user.Email)
	if err != nil {
		h.log.Error("failed to generate verification token", "email", user.Email, "error", err)
		c.JSON(500, gin.H{"error": "Failed to generate verification token"})
		return
	}

	err = h.emailService.SendVerificationEmail(ctx, user.Email, token)
	if err != nil {
		h.log.Error("failed to send verification email", "email", user.Email, "error", err)
		c.JSON(500, gin.H{"error": "User created but failed to send verification email"})
		return
	}

	h.log.Info("user registration completed successfully", "email", user.Email)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling user login request")

	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid login request", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.ValidateCredentials(ctx, req.Email, req.Password)
	if err != nil {
		h.log.Warn("credential validation failed", "email", req.Email, "error", err)
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.tokenService.GenerateTokenPair(ctx, user)
	if err != nil {
		h.log.Error("failed to generate token pair", "email", req.Email, "user_id", user.ID, "error", err)
		c.JSON(500, gin.H{"error": "Failed to generate authentication tokens"})
		return
	}

	c.SetCookie(
		"refresh_token",
		tokenPair.RefreshToken,
		int(tokenPair.RefreshExpiresIn),
		"/",
		"",
		false,
		true,
	)

	h.log.Info("user login successful", "email", req.Email, "user_id", user.ID)
	c.JSON(http.StatusOK, response.Auth{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling user logout request")

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.log.Warn("no refresh token found in cookie")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token found"})
		return
	}

	if err := h.tokenService.RevokeToken(ctx, refreshToken); err != nil {
		h.log.Error("failed to revoke refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	h.log.Info("user logged out successfully")
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling token refresh request")

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.log.Warn("no refresh token found in cookie")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token found"})
		return
	}

	tokenPair, err := h.tokenService.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		h.log.Warn("token refresh failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	h.log.Info("token refreshed successfully")
	c.JSON(http.StatusOK, response.Auth{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling email verification request")

	token := c.Query("token")
	if token == "" {
		h.log.Warn("missing token parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token parameter is required"})
		return
	}

	email, err := h.tokenService.ValidateVerificationToken(ctx, token)
	if err != nil {
		h.log.Warn("invalid or expired verification token", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	if err := h.userService.MarkEmailVerified(ctx, email); err != nil {
		h.log.Error("failed to mark email as verified", "email", email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	if err := h.tokenService.DeleteVerificationToken(ctx, token); err != nil {
		h.log.Warn("failed to delete verification token", "email", email, "error", err)
	}

	h.log.Info("email verified successfully", "email", email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}
