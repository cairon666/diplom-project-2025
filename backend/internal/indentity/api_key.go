package indentity

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type RolesService interface {
	GetPermissionsByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Permission, error)
	GetRolesByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Role, error)
}

type ExternalAppsService interface {
	GetByAPIKeyHash(ctx context.Context, apiKeyHash string) (models.ExternalApp, error)
}

type ApiKeyService struct {
	rolesService       RolesService
	externalAppService ExternalAppsService
}

func NewApiKeyService(rolesService RolesService, externalAppService ExternalAppsService) *ApiKeyService {
	return &ApiKeyService{
		rolesService:       rolesService,
		externalAppService: externalAppService,
	}
}

func (aks *ApiKeyService) GetAuthClaimsByAPIKeyHash(apiKeyHash string) (*AuthClaims, error) {
	apiKey, err := aks.externalAppService.GetByAPIKeyHash(context.Background(), apiKeyHash)
	if err != nil {
		return nil, err
	}

	perms, err := aks.rolesService.GetPermissionsByExternalAppID(context.Background(), apiKey.ID)
	if err != nil {
		return nil, err
	}

	roles, err := aks.rolesService.GetRolesByExternalAppID(context.Background(), apiKey.ID)
	if err != nil {
		return nil, err
	}

	return &AuthClaims{
		UserID:      apiKey.OwnerID,
		Permissions: perms,
		Roles:       roles,
	}, nil
}
