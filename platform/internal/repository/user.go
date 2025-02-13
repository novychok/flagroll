package repository

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type User interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, userID entity.UserID, user *entity.User) error
	Delete(ctx context.Context, userID entity.UserID) error
	GetByID(ctx context.Context, userID entity.UserID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
}
