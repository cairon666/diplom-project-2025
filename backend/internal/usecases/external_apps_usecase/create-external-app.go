package external_apps_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/models/permission"
	"github.com/cairon666/vkr-backend/pkg/api_key"
	"github.com/cairon666/vkr-backend/pkg/logger"
	"github.com/google/uuid"
)

var (
	availableRoleToExternalApp = []string{models.RoleExternalAppReader, models.RoleExternalAppWriter}
)

type CreateExternalAppRequest struct {
	Name  string
	Roles []string
}

func NewCreateExternalAppRequest(name string, roles []string) CreateExternalAppRequest {
	return CreateExternalAppRequest{
		Name:  name,
		Roles: roles,
	}
}

type CreateExternalAppResponse struct {
	ApiKey        string
	IdExternalApp uuid.UUID
}

func (uc *ExternalAppsUsecase) CreateExternalApp(ctx context.Context, dto CreateExternalAppRequest) (CreateExternalAppResponse, error) {
	claims, ok := indentity.GetAuthClaims(ctx)
	fmt.Println(claims)
	if !ok || !claims.HasPermission(permission.UpdateOwnExternalApps) {
		return CreateExternalAppResponse{}, apperrors.Forbidden()
	}

	// проверка что все роли можно применить
	for _, role := range dto.Roles {
		found := false
		for _, availableRole := range availableRoleToExternalApp {
			if role == availableRole {
				found = true

				break
			}
		}

		if !found {
			uc.logger.Error("try assign role from non available role", logger.String("role", role))

			return CreateExternalAppResponse{}, apperrors.Forbiddenf("try assign role from non available role: %s", role)
		}
	}

	id := uuid.New()
	apiKey, err := api_key.GenerateAPIKey()
	if err != nil {
		uc.logger.Error("failed to generate api key", logger.Error(err))

		return CreateExternalAppResponse{}, apperrors.InternalError()
	}
	hashApiKey := api_key.HashAPIKey(apiKey)

	externalApp := models.NewExternalApp(id, claims.UserID, dto.Name, hashApiKey, time.Now())
	if err := uc.externalAppsService.Create(ctx, externalApp); err != nil {
		uc.logger.Error("failed to create external app", logger.Error(err))

		return CreateExternalAppResponse{}, apperrors.InternalError()
	}

	if err := uc.rolesService.AssignRolesToExternalApp(ctx, externalApp.ID, dto.Roles); err != nil {
		uc.logger.Error("failed to assign roles to external app", logger.Error(err))

		return CreateExternalAppResponse{}, apperrors.InternalError()
	}

	return CreateExternalAppResponse{ApiKey: apiKey, IdExternalApp: id}, nil
}
