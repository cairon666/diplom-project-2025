package external_apps_service

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type ExternalAppsRepo interface {
	Create(ctx context.Context, app models.ExternalApp) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (models.ExternalApp, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.ExternalApp, error)
	UpdateName(ctx context.Context, id uuid.UUID, newName string) error
	GetByAPIKeyHash(ctx context.Context, hash string) (models.ExternalApp, error)
}

type ExternalAppsService struct {
	externalAppsRepo ExternalAppsRepo
}

func NewExternalAppsService(externalAppsRepo ExternalAppsRepo) *ExternalAppsService {
	return &ExternalAppsService{
		externalAppsRepo: externalAppsRepo,
	}
}

func (s *ExternalAppsService) Create(ctx context.Context, app models.ExternalApp) error {
	return s.externalAppsRepo.Create(ctx, app)
}

func (s *ExternalAppsService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.externalAppsRepo.Delete(ctx, id)
}

func (s *ExternalAppsService) GetByID(ctx context.Context, id uuid.UUID) (models.ExternalApp, error) {
	return s.externalAppsRepo.GetByID(ctx, id)
}

func (s *ExternalAppsService) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.ExternalApp, error) {
	return s.externalAppsRepo.ListByOwner(ctx, ownerID)
}

func (s *ExternalAppsService) UpdateName(ctx context.Context, id uuid.UUID, newName string) error {
	return s.externalAppsRepo.UpdateName(ctx, id, newName)
}

func (s *ExternalAppsService) GetByAPIKeyHash(ctx context.Context, hash string) (models.ExternalApp, error) {
	return s.externalAppsRepo.GetByAPIKeyHash(ctx, hash)
}
