package user_passwords_repo

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserPasswordsRepo struct {
	dbQueries *dbqueries.Queries
}

func NewUserPasswordsRepo(dbQueries *dbqueries.Queries) *UserPasswordsRepo {
	return &UserPasswordsRepo{
		dbQueries: dbQueries,
	}
}

func (u *UserPasswordsRepo) GetUserPasswordByUserId(ctx context.Context, userId uuid.UUID) (models.UserPassword, error) {
	dbUserPassword, err := u.dbQueries.GetUserPassword(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserPassword{}, apperrors.NotFound()
		}

		return models.UserPassword{}, err
	}

	return models.NewUserPassword(dbUserPassword.UserID, dbUserPassword.Salt, dbUserPassword.PasswordHash), nil
}

func (u *UserPasswordsRepo) CreateUserPassword(ctx context.Context, userPassword models.UserPassword) error {
	err := u.dbQueries.CreateUserPassword(ctx, dbqueries.CreateUserPasswordParams{
		UserID:       userPassword.UserId,
		Salt:         userPassword.Salt,
		PasswordHash: userPassword.PasswordHash,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "USER_PASSWORDS_USER_ID_UNIQUE":
				return apperrors.AlreadyExists()
			default:
				return err
			}
		}

		return err
	}

	return nil
}

func (u *UserPasswordsRepo) UpdateUserPassword(ctx context.Context, userPassword models.UserPassword) error {
	err := u.dbQueries.UpdateUserPassword(ctx, dbqueries.UpdateUserPasswordParams{
		UserID:       userPassword.UserId,
		Salt:         userPassword.Salt,
		PasswordHash: userPassword.PasswordHash,
	})

	if err != nil {
		return err
	}

	return nil
}
