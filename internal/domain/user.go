package domain

import "context"

type User struct {
	ID            int64  `json:"id"`
	Login         string `json:"login"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	FullName      string `json:"full_name"`
	Password      string `json:"-"`
	Rating        int    `json:"rating"`
	Avatar        string `json:"avatar"`
	EmailVerified bool   `json:"email_verified"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
