package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo               domain.UserRepository
	tokenService       domain.TokenService
	emailSenderService domain.EmailSender
	log                *logger.Logger
}

func NewUserService(repo domain.UserRepository, tokenService domain.TokenService, EmailSenderService domain.EmailSender, log *logger.Logger) *UserService {
	return &UserService{
		repo:               repo,
		tokenService:       tokenService,
		emailSenderService: EmailSenderService,
		log:                log,
	}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	s.log.Info("attempting user registration", "email", user.Email, "login", user.Login)

	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		s.log.Error("failed to check existing user", "email", user.Email, "error", err)
		return fmt.Errorf("database error: %v", err)
	}
	if existing != nil {
		s.log.Warn("registration failed: email already exists", "email", user.Email)
		return errors.New("email already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return errors.New("failed to hash password")
	}
	user.Password = string(hashed)

	if err := s.repo.Create(ctx, user); err != nil {
		s.log.Error("failed to create user in database", "email", user.Email, "error", err)
		return fmt.Errorf("database error: %v", err)
	}

	s.log.Info("user created successfully", "email", user.Email, "user_id", user.ID)

	token, err := s.tokenService.GenerateVerificationToken(ctx, user.Email)
	if err != nil {
		s.log.Error("failed to generate verification token", "email", user.Email, "error", err)
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	if err := s.emailSenderService.SendVerificationEmail(ctx, user.Email, token); err != nil {
		s.log.Error("failed to send verification email", "email", user.Email, "error", err)
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	s.log.Info("registration completed successfully, verification email sent", "email", user.Email)
	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	s.log.Info("attempting user login", "email", email)

	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user for login", "email", email, "error", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	if existing == nil {
		s.log.Warn("login failed: user not found", "email", email)
		return nil, errors.New("invalid credentials")
	}

	if existing.EmailVerified == false {
		s.log.Warn("login failed: email not verified", "email", email, "user_id", existing.ID)
		return nil, errors.New("email not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(password))
	if err != nil {
		s.log.Warn("login failed: invalid password", "email", email, "user_id", existing.ID)
		return nil, errors.New("invalid credentials")
	}

	tokenPair, err := s.tokenService.GenerateTokenPair(ctx, existing)
	if err != nil {
		s.log.Error("failed to generate token pair", "email", email, "user_id", existing.ID, "error", err)
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	s.log.Info("user logged in successfully", "email", email, "user_id", existing.ID)
	return tokenPair, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) error {
	s.log.Info("attempting user logout")

	err := s.tokenService.RevokeToken(ctx, refreshToken)
	if err != nil {
		s.log.Error("logout failed: unable to revoke refresh token", "error", err)
		return err
	}

	s.log.Info("user logged out successfully")
	return nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	s.log.Info("attempting token refresh")

	tokenPair, err := s.tokenService.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		s.log.Warn("token refresh failed", "error", err)
		return nil, err
	}

	s.log.Info("token refreshed successfully")
	return tokenPair, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	s.log.Info("attempting email verification")

	email, err := s.tokenService.ValidateVerificationToken(ctx, token)
	if err != nil {
		s.log.Warn("email verification failed: invalid or expired token", "error", err)
		return fmt.Errorf("invalid or expired verification token: %v", err)
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user for email verification", "email", email, "error", err)
		return fmt.Errorf("database error: %v", err)
	}
	if user == nil {
		s.log.Warn("email verification failed: user not found", "email", email)
		return errors.New("user not found")
	}

	if user.EmailVerified {
		s.log.Info("email already verified", "email", email, "user_id", user.ID)
		return nil
	}

	user.EmailVerified = true
	if err := s.repo.Update(ctx, user); err != nil {
		s.log.Error("failed to update user email verification status", "email", email, "user_id", user.ID, "error", err)
		return fmt.Errorf("failed to update user: %v", err)
	}

	s.log.Info("email verified successfully", "email", email, "user_id", user.ID)

	if err := s.tokenService.DeleteVerificationToken(ctx, token); err != nil {
		s.log.Warn("failed to delete verification token", "email", email, "error", err)
	}

	return nil
}
