package external_apps_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/google/uuid"
)

type ExternalAppListItem struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	Name      string
	CreatedAt time.Time
	Roles     []string
}

type GetExternalAppListResponse struct {
	ExternalApps []ExternalAppListItem
}

func (eau *ExternalAppsUsecase) GetExternalAppList(ctx context.Context) (GetExternalAppListResponse, error) {
	claims, ok := indentity.GetAuthClaims(ctx)
	if !ok || !claims.HasPermission(permission.ReadOwnExternalApps) {
		return GetExternalAppListResponse{}, apperrors.Forbidden()
	}

	externalAppItems := make([]ExternalAppListItem, 0)
	externalApps, err := eau.externalAppsService.ListByOwner(ctx, claims.UserID)
	if errors.Is(err, apperrors.NotFound()) {
		return GetExternalAppListResponse{ExternalApps: externalAppItems}, nil
	} else if err != nil {
		return GetExternalAppListResponse{}, err
	}

	for _, externalApp := range externalApps {
		roles, err := eau.rolesService.GetRolesByExternalAppID(ctx, externalApp.ID)
		var roleNames []string
		if err != nil {
			// Если ошибка получения ролей, то считаем что ролей нет
			roleNames = []string{}
		} else {
			roleNames = make([]string, len(roles))
			for i, role := range roles {
				roleNames[i] = role.Name
			}
		}

		externalAppItems = append(externalAppItems, ExternalAppListItem{
			ID:        externalApp.ID,
			OwnerID:   externalApp.OwnerID,
			Name:      externalApp.Name,
			CreatedAt: externalApp.CreatedAt,
			Roles:     roleNames,
		})
	}

	return GetExternalAppListResponse{ExternalApps: externalAppItems}, nil
}
