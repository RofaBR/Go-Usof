package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	setter := &models.CategorySetter{
		Title:       omit.From(category.Title),
		Slug:        omit.From(category.Slug),
		Description: omitnull.From(category.Desc),
	}

	query := models.Categories.Insert(setter)

	model, err := query.One(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}
	category.ID = int64(model.ID)

	return nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]*domain.Category, error) {
	query := models.Categories.Query()

	catSlice, err := query.All(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if len(catSlice) == 0 {
		return []*domain.Category{}, nil
	}
	categories := make([]*domain.Category, len(catSlice))
	for i, category := range catSlice {
		categories[i] = &domain.Category{
			ID:    int64(category.ID),
			Title: category.Title,
			Slug:  category.Slug,
			Desc:  category.Description.GetOrZero(),
		}
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := models.Categories.Query(
		sm.Where(models.Categories.Columns.ID.EQ(psql.Arg(id))),
	)
	category, err := query.One(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &domain.Category{
		ID:    int64(category.ID),
		Title: category.Title,
		Slug:  category.Slug,
		Desc:  category.Description.GetOrZero(),
	}, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := models.Categories.Query(
		sm.Where(models.Categories.Columns.Slug.EQ(psql.Arg(slug))),
	)
	category, err := query.One(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &domain.Category{
		ID:    int64(category.ID),
		Title: category.Title,
		Slug:  category.Slug,
		Desc:  category.Description.GetOrZero(),
	}, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	setter := &models.CategorySetter{
		Title:       omit.From(category.Title),
		Slug:        omit.From(category.Slug),
		Description: omitnull.From(category.Desc),
	}
	query := models.Categories.Update(
		setter.UpdateMod(),
		um.Where(models.Categories.Columns.ID.EQ(psql.Arg(category.ID))),
	)
	rowsAffected, err := query.Exec(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", category.ID)
	}

	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	query := models.Categories.Delete(
		dm.Where(models.Categories.Columns.ID.EQ(psql.Arg(id))),
	)

	rowsAffected, err := query.Exec(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", id)
	}
	return nil
}
