package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/models"
	"github.com/RofaBR/Go-Usof/internal/models/enums"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/um"

	"github.com/aarondl/opt/omit"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	setter := &models.UserSetter{
		Login:    omit.From(user.Login),
		Email:    omit.From(user.Email),
		Password: omit.From(user.Password),
		Fullname: omit.From(user.FullName),
		Role:     omit.From(enums.UserRole(user.Role)),
		Rating:   omit.From(int32(user.Rating)),
	}

	query := models.Users.Insert(setter)

	model, err := query.One(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}
	user.ID = int(model.ID)

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := models.Users.Query(
		sm.Where(models.Users.Columns.Email.EQ(psql.Arg(email))),
	)

	model, err := query.One(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return mapModelToDomain(model), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := models.Users.Query(
		sm.Where(models.Users.Columns.ID.EQ(psql.Arg(id))),
	)

	model, err := query.One(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return mapModelToDomain(model), nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := models.Users.Query()

	userSlice, err := query.All(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if len(userSlice) == 0 {
		return []*domain.User{}, nil
	}
	users := make([]*domain.User, len(userSlice))
	for i, model := range userSlice {
		users[i] = mapModelToDomain(model)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	setter := &models.UserSetter{
		Email:         omit.From(user.Email),
		Password:      omit.From(user.Password),
		Fullname:      omit.From(user.FullName),
		EmailVerified: omit.From(user.EmailVerified),
	}

	query := models.Users.Update(
		setter.UpdateMod(),
		um.Where(models.Users.Columns.ID.EQ(psql.Arg(user.ID))),
	)
	rowsAffected, err := query.Exec(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := models.Users.Delete(
		dm.Where(models.Users.Columns.ID.EQ(psql.Arg(id))),
	)

	rowsAffected, err := query.Exec(ctx, bob.NewDB(stdlib.OpenDBFromPool(r.db)))
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}
	return nil
}

func mapModelToDomain(m *models.User) *domain.User {
	return &domain.User{
		ID:            int(m.ID),
		Login:         m.Login,
		Email:         m.Email,
		Password:      m.Password,
		FullName:      m.Fullname,
		Role:          string(m.Role),
		Rating:        int(m.Rating),
		EmailVerified: m.EmailVerified,
	}
}
