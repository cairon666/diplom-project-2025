package roles_repo

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetRolesByUserID получает список ролей пользователя по его UUID.
func (r *RolesRepo) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	roles, err := r.query.GetRolesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, apperrors.InternalErrorWithCause(err)
	}
	result := make([]models.Role, len(roles))
	for i, r := range roles {
		result[i] = models.Role{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	return result, nil
}

func (r *RolesRepo) GetRoleByName(ctx context.Context, name string) (models.Role, error) {
	role, err := r.query.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Role{}, apperrors.NotFound()
		}

		return models.Role{}, apperrors.InternalErrorWithCause(err)
	}

	return models.Role{
		ID:   role.ID,
		Name: role.Name,
	}, nil
}

func (r *RolesRepo) GetRolesByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Role, error) {
	roles, err := r.query.GetRolesByExternalAppID(ctx, externalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, apperrors.InternalErrorWithCause(err)
	}

	result := make([]models.Role, len(roles))
	for i, r := range roles {
		result[i] = models.Role{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	return result, nil
}
