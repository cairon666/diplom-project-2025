package external_apps_repo

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ExternalAppsRepo struct {
	query *dbqueries.Queries
}

func NewExternalAppsRepo(query *dbqueries.Queries) *ExternalAppsRepo {
	return &ExternalAppsRepo{
		query: query,
	}
}

func (r *ExternalAppsRepo) Create(ctx context.Context, app models.ExternalApp) error {
	_, err := r.query.CreateExternalApp(ctx, dbqueries.CreateExternalAppParams{
		ID:          app.ID,
		Name:        app.Name,
		OwnerUserID: app.OwnerID,
		ApiKeyHash:  app.APIKeyHash,
		CreatedAt:   app.CreatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ExternalAppsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.query.DeleteExternalApp(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *ExternalAppsRepo) GetByID(ctx context.Context, id uuid.UUID) (models.ExternalApp, error) {
	dbApp, err := r.query.GetExternalAppByID(ctx, id)
	if err != nil {
		return models.ExternalApp{}, err
	}
	return toModel(dbApp), nil
}

func (r *ExternalAppsRepo) GetByAPIKeyHash(ctx context.Context, hash string) (models.ExternalApp, error) {
	dbApp, err := r.query.GetExternalAppByAPIKeyHash(ctx, hash)
	if err != nil {
		return models.ExternalApp{}, err
	}
	return toModel(dbApp), nil
}

func (r *ExternalAppsRepo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.ExternalApp, error) {
	dbApps, err := r.query.ListExternalAppsByOwner(ctx, ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	result := make([]models.ExternalApp, 0, len(dbApps))
	for _, app := range dbApps {
		result = append(result, toModel(app))
	}
	return result, nil
}

func (r *ExternalAppsRepo) UpdateName(ctx context.Context, id uuid.UUID, newName string) error {
	return r.query.UpdateExternalAppName(ctx, dbqueries.UpdateExternalAppNameParams{
		ID:   id,
		Name: newName,
	})
}

// Маппер из sqlc-структуры в модель
func toModel(dbApp dbqueries.EXTERNALAPP) models.ExternalApp {
	return models.ExternalApp{
		ID:         dbApp.ID,
		OwnerID:    dbApp.OwnerUserID,
		Name:       dbApp.Name,
		APIKeyHash: dbApp.ApiKeyHash,
		CreatedAt:  dbApp.CreatedAt,
	}
}
