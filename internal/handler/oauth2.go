package handler

import (
	"net/http"

	"github.com/RofaBR/Go-Usof/internal/dto/response"
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
)

type OAuth2Handler struct {
	oauth2Service *services.OAuth2Service
	tokenService  *services.TokenService
	log           *logger.Logger
}

func NewOAuth2Handler(oauth2Service *services.OAuth2Service, tokenService *services.TokenService, log *logger.Logger) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2Service: oauth2Service,
		tokenService:  tokenService,
		log:           log,
	}
}

func (h *OAuth2Handler) GoogleLogin(c *gin.Context) {
	url, state, err := h.oauth2Service.GetAuthURL(c.Request.Context())
	if err != nil {
		h.log.Error("failed to generate auth url", "error", err)
		c.JSON(500, gin.H{"error": "Failed to generate auth url"})
		return
	}

	c.SetCookie(
		"oauth_state",
		state,
		300,
		"/",
		"",
		false,
		true,
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *OAuth2Handler) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	storedState, err := c.Cookie("oauth_state")
	if err != nil {
		h.log.Error("failed to get oauth state cookie", "error", err)
		c.JSON(500, gin.H{"error": "Failed to get oauth state cookie"})
		return
	}

	urlState := c.Query("state")
	if urlState != storedState {
		h.log.Error("oauth state mismatch", "urlState", urlState, "storedState", storedState)
		c.JSON(400, gin.H{"error": "OAuth state mismatch"})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	code := c.Query("code")
	if code == "" {
		h.log.Warn("no code found in query params")
		c.JSON(400, gin.H{"error": "No code found in query params"})
		return
	}

	user, _, err := h.oauth2Service.HandleCallback(ctx, code)
	if err != nil {
		h.log.Error("failed to handle oauth callback", "error", err)
		c.JSON(500, gin.H{"error": "Failed to handle oauth callback"})
		return
	}
	tokenPair, err := h.tokenService.GenerateTokenPair(ctx, user)
	if err != nil {
		h.log.Error("failed to generate token pair", "error", err)
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
	c.JSON(http.StatusOK, response.Auth{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
	})
}
