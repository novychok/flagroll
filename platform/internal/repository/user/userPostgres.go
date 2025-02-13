package user

import (
	"context"
	"database/sql"

	"github.com/novychok/flagroll/platform/internal/database/dao"
	"github.com/novychok/flagroll/platform/internal/database/pqmodels"
	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/pkg/postgres"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Create(ctx context.Context, user *entity.User) error {
	userDB := &pqmodels.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	err := userDB.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}

	user.ID = entity.UserID(userDB.ID)
	user.CreatedAt = userDB.CreatedAt
	user.UpdatedAt = userDB.UpdatedAt

	return nil
}

func (r *postgresRepository) Update(ctx context.Context, userID entity.UserID, user *entity.User) error {

	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, userID entity.UserID) error {

	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, userID entity.UserID) (*entity.User, error) {
	userDB, err := pqmodels.Users(
		pqmodels.UserWhere.ID.EQ(userID.String()),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &entity.User{}
	dao.UserTo(userDB, user)

	return user, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	userDB, err := pqmodels.Users(
		pqmodels.UserWhere.Email.EQ(email),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &entity.User{}
	dao.UserTo(userDB, user)

	return user, nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*entity.User, error) {

	return nil, nil
}

func NewPostgres(db postgres.Connection) repository.User {
	return &postgresRepository{
		db: db,
	}
}
