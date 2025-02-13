package platformapiv1

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type UserCtxKey struct{}

func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, UserCtxKey{}, user)
}

func UserFromContext(ctx context.Context) *entity.User {
	user, _ := ctx.Value(UserCtxKey{}).(*entity.User)
	return user
}
