package services

import (
	"context"
	"fmt"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("check email failed: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("email already registered")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}
	user.Password = string(hashed)

	return s.repo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, login, email, password string) (*domain.User, error) {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email failed: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return existing, nil
}
