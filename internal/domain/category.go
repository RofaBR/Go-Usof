package domain

import "context"

type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Desc  string `json:"description"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetAll(ctx context.Context) ([]*Category, error)
	GetByID(ctx context.Context, id int64) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int64) error
}
