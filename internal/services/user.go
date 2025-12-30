package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo               domain.UserRepository
	tokenService       domain.TokenService
	emailSenderService domain.EmailSender
}

func NewUserService(repo domain.UserRepository, tokenService domain.TokenService, EmailSenderService domain.EmailSender) *UserService {
	return &UserService{
		repo:               repo,
		tokenService:       tokenService,
		emailSenderService: EmailSenderService,
	}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	if existing != nil {
		return errors.New("email already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashed)

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	token, err := s.tokenService.GenerateVerificationToken(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	if err := s.emailSenderService.SendVerificationEmail(ctx, user.Email, token); err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	if existing == nil {
		return nil, errors.New("invalid credentials")
	}

	if existing.EmailVerified == false {
		return nil, errors.New("email not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	tokenPair, err := s.tokenService.GenerateTokenPair(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return tokenPair, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenService.RevokeToken(ctx, refreshToken)
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	return s.tokenService.RefreshAccessToken(ctx, refreshToken)
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	email, err := s.tokenService.ValidateVerificationToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token: %v", err)
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	if err := s.tokenService.DeleteVerificationToken(ctx, token); err != nil {
		fmt.Printf("Warning: failed to delete verification token: %v\n", err)
	}

	return nil
}
