package auth_providers_repo

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

type AuthProvidersRepo struct {
	dbQuery *dbqueries.Queries
}

func NewAuthProvidersRepo(dbQuery *dbqueries.Queries) *AuthProvidersRepo {
	return &AuthProvidersRepo{
		dbQuery: dbQuery,
	}
}

func (a *AuthProvidersRepo) GetAuthProviderById(ctx context.Context, id uuid.UUID) (models.AuthProvider, error) {
	dbAuthProvider, err := a.dbQuery.GetAuthProviderById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.AuthProvider{}, apperrors.NotFound()
		}
		return models.AuthProvider{}, err
	}

	return models.NewAuthProvider(
		dbAuthProvider.ID,
		dbAuthProvider.UserID,
		dbAuthProvider.ProviderName,
		dbAuthProvider.ProviderUserID,
		dbAuthProvider.CreatedAt,
	), nil
}

func (a *AuthProvidersRepo) GetAuthProvidersByUserId(ctx context.Context, userId uuid.UUID) ([]models.AuthProvider, error) {
	dbAuthProviders, err := a.dbQuery.GetAuthProvidersByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}
		return nil, err
	}

	if len(dbAuthProviders) == 0 {
		return nil, apperrors.NotFound()
	}

	authProviders := make([]models.AuthProvider, 0, len(dbAuthProviders))
	for _, dbAuthProvider := range dbAuthProviders {
		authProviders = append(authProviders,
			models.NewAuthProvider(
				dbAuthProvider.ID,
				dbAuthProvider.UserID,
				dbAuthProvider.ProviderName,
				dbAuthProvider.ProviderUserID,
				dbAuthProvider.CreatedAt,
			))
	}

	return authProviders, nil
}

func (a *AuthProvidersRepo) GetAuthProviderByUserIdAndProviderName(ctx context.Context, userId uuid.UUID, providerName string) (models.AuthProvider, error) {
	dbAuthProvider, err := a.dbQuery.GetAuthProviderByUserIdAndProviderName(ctx, dbqueries.GetAuthProviderByUserIdAndProviderNameParams{
		UserID:       userId,
		ProviderName: providerName,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.AuthProvider{}, apperrors.NotFound()
		}
		return models.AuthProvider{}, err
	}

	return models.NewAuthProvider(
		dbAuthProvider.ID,
		dbAuthProvider.UserID,
		dbAuthProvider.ProviderName,
		dbAuthProvider.ProviderUserID,
		dbAuthProvider.CreatedAt,
	), nil
}

func (a *AuthProvidersRepo) GetAuthProviderByProviderUserIdAndProviderName(ctx context.Context, providerUserID int64, providerName string) (models.AuthProvider, error) {
	dbAuthProvider, err := a.dbQuery.GetAuthProviderByProviderUserIdAndProviderName(ctx, dbqueries.GetAuthProviderByProviderUserIdAndProviderNameParams{
		ProviderUserID: providerUserID,
		ProviderName:   providerName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.AuthProvider{}, apperrors.NotFound()
		}
		return models.AuthProvider{}, err
	}

	return models.NewAuthProvider(
		dbAuthProvider.ID,
		dbAuthProvider.UserID,
		dbAuthProvider.ProviderName,
		dbAuthProvider.ProviderUserID,
		dbAuthProvider.CreatedAt,
	), nil
}

func (a *AuthProvidersRepo) CreateAuthProvider(ctx context.Context, authProvider models.AuthProvider) error {
	err := a.dbQuery.CreateAuthProvider(ctx, dbqueries.CreateAuthProviderParams{
		ID:             authProvider.ID,
		UserID:         authProvider.UserId,
		ProviderName:   authProvider.ProviderName,
		ProviderUserID: authProvider.ProviderUserId,
		CreatedAt:      authProvider.CreatedAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "AUTH_PROVIDERS_PROVIDER_USER_ID_UNIQUE":
				return apperrors.ProviderAccountAlreadyLinked()
			case pgErr.ConstraintName == "AUTH_PROVIDERS_USER_ID_PROVIDER_NAME_UNIQUE":
				return apperrors.ProviderAlreadyConnected()
			default:
				return err
			}
		}
		return err
	}

	return nil
}

func (a *AuthProvidersRepo) DeleteAuthProviderById(ctx context.Context, id uuid.UUID) error {
	err := a.dbQuery.DeleteAuthProviderById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthProvidersRepo) DeleteAuthProviderByUserId(ctx context.Context, userId uuid.UUID) error {
	err := a.dbQuery.DeleteAuthProviderByUserId(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
