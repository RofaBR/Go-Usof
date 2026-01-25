package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/RofaBR/Go-Usof/internal/config"
	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type OAuth2Service struct {
	oauthConfig *oauth2.Config
	userRepo    domain.UserRepository
	logger      *logger.Logger
}

func NewOAuth2Service(cfg *config.OAuth2Config, userRepo domain.UserRepository, log *logger.Logger) *OAuth2Service {
	return &OAuth2Service{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURI,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		userRepo: userRepo,
		logger:   log,
	}
}

func (s *OAuth2Service) GetAuthURL(ctx context.Context) (string, string, error) {
	state, err := generateState()
	if err != nil {
		s.logger.Error("failed to generate state", "error", err)
	}
	url := s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url, state, nil
}

func (s *OAuth2Service) HandleCallback(ctx context.Context, code string) (*domain.User, bool, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, false, fmt.Errorf("failed to exchange token: %w", err)
	}
	userInfo, err := s.fetchGoogleUser(ctx, token)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch user info: %w", err)
	}
	user, err := s.userRepo.GetByGoogleID(ctx, userInfo.ID)
	if err != nil {
		return nil, false, fmt.Errorf("database error: %w", err)
	}
	if user != nil {
		return user, false, nil
	}
	user, err = s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, false, fmt.Errorf("database error: %w", err)
	}
	if user != nil {
		user.GoogleID = userInfo.ID
		s.userRepo.Update(ctx, user)
		return user, false, nil
	}
	randomPass, err := generateRandomPassword()
	if err != nil {
		return nil, false, fmt.Errorf("failed to generate password: %w", err)
	}
	newUser := &domain.User{
		Login:         userInfo.Email,
		Email:         userInfo.Email,
		GoogleID:      userInfo.ID,
		Password:      randomPass,
		FullName:      userInfo.Name,
		Role:          "user",
		Avatar:        userInfo.Picture,
		EmailVerified: userInfo.VerifiedEmail,
	}
	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, false, err
	}
	return newUser, true, nil
}

func (s *OAuth2Service) fetchGoogleUser(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error("failed to close response body", "error", err)
		}
	}(resp.Body)
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}
	return &userInfo, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateRandomPassword() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
