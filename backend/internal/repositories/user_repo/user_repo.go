package user_repo

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

type UserRepo struct {
	query *dbqueries.Queries
}

func NewUserRepo(query *dbqueries.Queries) *UserRepo {
	return &UserRepo{
		query: query,
	}
}

func getEmail(email *string) string {
	if email == nil {
		return ""
	}

	return *email
}

func toDbEmail(email string) *string {
	if email == "" {
		return nil
	}

	return &email
}

func (u *UserRepo) GetUserById(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	dbUser, err := u.query.GetUserById(ctx, ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, err
	}

	user := models.NewUser(dbUser.ID, getEmail(dbUser.Email), dbUser.FirstName, dbUser.LastName, dbUser.IsRegistrationComplete, dbUser.CreatedAt)

	return user, nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	dbUser, err := u.query.GetUserByEmail(ctx, &email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, err
	}

	user := models.NewUser(dbUser.ID, getEmail(dbUser.Email), dbUser.FirstName, dbUser.LastName, dbUser.IsRegistrationComplete, dbUser.CreatedAt)

	return user, nil
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	params := dbqueries.CreateUserParams{
		ID:                     user.ID,
		Email:                  toDbEmail(user.Email),
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		IsRegistrationComplete: user.IsRegistrationComplete,
	}

	err := u.query.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "USERS_EMAIL_UNIQUE":
				return apperrors.AlreadyExists()
			default:
				return err
			}
		}

		return err
	}

	return nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	err := u.query.UpdateUserFull(ctx, dbqueries.UpdateUserFullParams{
		ID:                     user.ID,
		Email:                  toDbEmail(user.Email),
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		IsRegistrationComplete: user.IsRegistrationComplete,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "USERS_EMAIL_UNIQUE":
				return apperrors.AlreadyExists()
			}
		}

		return err
	}

	return nil
}
