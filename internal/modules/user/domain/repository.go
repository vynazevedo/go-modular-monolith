package domain

import (
	"context"
)

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, page, limit int) ([]*User, error)
	Delete(ctx context.Context, id string) error
}
