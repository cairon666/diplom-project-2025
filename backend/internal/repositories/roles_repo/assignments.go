package roles_repo

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/cairon666/vkr-backend/pkg/postgres"
	"github.com/google/uuid"
)

// AssignRoleToUser присваивает роль пользователю.
func (r *RolesRepo) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int32) error {
	err := r.query.AssignRoleToUser(ctx, dbqueries.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return apperrors.InternalErrorWithCause(err)
	}

	return nil
}

// AssignRoleToExternalApp присваивает роль внешнему приложению.
func (r *RolesRepo) AssignRoleToExternalApp(ctx context.Context, externalID uuid.UUID, roleID int32) error {
	err := r.query.AddRoleToExternalApp(ctx, dbqueries.AddRoleToExternalAppParams{
		ExternalAppID: externalID,
		RoleID:        roleID,
	})
	if err != nil {
		return apperrors.InternalErrorWithCause(err)
	}

	return nil
}

// AssignRolesToExternalApp присваивает роли внешнему приложению.
func (r *RolesRepo) AssignRolesToExternalApp(ctx context.Context, externalID uuid.UUID, roleIDs []int32) error {
	builder := postgres.Builder.Insert("EXTERNAL_APPS_ROLES").Cols(
		"external_app_id",
		"role_id",
	)

	for _, roleID := range roleIDs {
		builder = builder.Vals([]interface{}{
			externalID,
			roleID,
		})
	}

	sql, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		return apperrors.InternalErrorWithCause(err)
	}

	if _, err := r.dbClient.Exec(ctx, sql, args...); err != nil {
		return apperrors.InternalErrorWithCause(err)
	}

	return nil
}
