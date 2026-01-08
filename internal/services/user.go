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
	repo domain.UserRepository
	log  *logger.Logger
}

func NewUserService(repo domain.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	s.log.Info("creating user", "email", user.Email, "login", user.Login)
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		s.log.Error("failed to check existing user", "email", user.Email, "error", err)
		return fmt.Errorf("database error: %v", err)
	}
	if existing != nil {
		s.log.Warn("user creation failed: email already exists", "email", user.Email)
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
	return nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	s.log.Info("getting user by email", "email", email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user by email", "email", email, "error", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	if user == nil {
		s.log.Info("user not found", "email", email)
		return nil, nil
	}

	s.log.Info("user retrieved successfully", "email", email, "user_id", user.ID)
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	s.log.Info("getting user by ID", "user_id", id)

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get user by ID", "user_id", id, "error", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	if user == nil {
		s.log.Info("user not found", "user_id", id)
		return nil, nil
	}

	s.log.Info("user retrieved successfully", "user_id", id)
	return user, nil
}

func (s *UserService) ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	s.log.Info("validating user credentials", "email", email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user for credential validation", "email", email, "error", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	if user == nil {
		s.log.Warn("credential validation failed: user not found", "email", email)
		return nil, errors.New("invalid credentials")
	}

	if !user.EmailVerified {
		s.log.Warn("credential validation failed: email not verified", "email", email, "user_id", user.ID)
		return nil, errors.New("email not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.log.Warn("credential validation failed: invalid password", "email", email, "user_id", user.ID)
		return nil, errors.New("invalid credentials")
	}

	s.log.Info("credentials validated successfully", "email", email, "user_id", user.ID)
	return user, nil
}

func (s *UserService) MarkEmailVerified(ctx context.Context, email string) error {
	s.log.Info("marking email as verified", "email", email)

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
	return nil
}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	s.log.Info("updating user", "user_id", user.ID)

	if err := s.repo.Update(ctx, user); err != nil {
		s.log.Error("failed to update user", "user_id", user.ID, "error", err)
		return fmt.Errorf("database error: %v", err)
	}

	s.log.Info("user updated successfully", "user_id", user.ID)
	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	s.log.Info("deleting user", "user_id", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete user", "user_id", id, "error", err)
		return fmt.Errorf("database error: %v", err)
	}

	s.log.Info("user deleted successfully", "user_id", id)
	return nil
}
