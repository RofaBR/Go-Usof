package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/pkg/logger"
)

type CategoryService struct {
	repo domain.CategoryRepository
	log  *logger.Logger
}

func NewCategoryService(repo domain.CategoryRepository, log *logger.Logger) *CategoryService {
	return &CategoryService{repo: repo, log: log}
}

func (s *CategoryService) Create(ctx context.Context, category *domain.Category) error {
	s.log.Info("create category")
	existing, err := s.repo.GetBySlug(ctx, category.Slug)
	if err != nil {
		s.log.Error("error getting category by slug", err)
		return fmt.Errorf("database error: %v", err)
	}
	if existing != nil {
		s.log.Warn("category already exists")
		return errors.New("category already exists")
	}
	if err := s.repo.Create(ctx, category); err != nil {
		s.log.Error("error creating category", err)
		return fmt.Errorf("database error: %v", err)
	}
	s.log.Info("category created")
	return nil
}
