package external_apps_usecase

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

type ExternalAppsService interface {
	Create(ctx context.Context, app models.ExternalApp) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (models.ExternalApp, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.ExternalApp, error)
	UpdateName(ctx context.Context, id uuid.UUID, newName string) error
}

type RolesService interface {
	AssignRolesToExternalApp(ctx context.Context, externalID uuid.UUID, roleNames []string) error
}

type ExternalAppsUsecase struct {
	logger              *logger.Logger
	externalAppsService ExternalAppsService
	rolesService        RolesService
}

func NewExternalAppsUsecase(logger *logger.Logger, externalAppsService ExternalAppsService, rolesService RolesService) *ExternalAppsUsecase {
	return &ExternalAppsUsecase{
		logger:              logger,
		externalAppsService: externalAppsService,
		rolesService:        rolesService,
	}
}
