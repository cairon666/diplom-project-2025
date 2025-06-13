package roles_repo

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *RolesRepo) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []int32) ([]models.Permission, error) {
	perms, err := r.query.GetPermissionsByRoleIDs(ctx, roleIDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, apperrors.InternalErrorWithCause(err)
	}
	result := make([]models.Permission, len(perms))
	for i, p := range perms {
		result[i] = models.Permission{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	return result, nil
}

// GetPermissionsByUserID получает список разрешений пользователя по его UUID.
func (r *RolesRepo) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Permission, error) {
	perms, err := r.query.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, apperrors.InternalErrorWithCause(err)
	}
	result := make([]models.Permission, len(perms))
	for i, p := range perms {
		result[i] = models.Permission{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	return result, nil
}

func (r *RolesRepo) GetPermissionsByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Permission, error) {
	perms, err := r.query.GetPermissionsByExternalAppID(ctx, externalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, apperrors.InternalErrorWithCause(err)
	}
	result := make([]models.Permission, len(perms))
	for i, p := range perms {
		result[i] = models.Permission{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	return result, nil
}

// HasUserPermission проверяет, есть ли у пользователя разрешение с заданным именем.
func (r *RolesRepo) HasUserPermission(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	hasUserPermission, err := r.query.HasUserPermission(ctx, dbqueries.HasUserPermissionParams{
		UserID: userID,
		Name:   permissionName,
	})
	if err != nil {
		return false, apperrors.InternalErrorWithCause(err)
	}

	return hasUserPermission, nil
}
