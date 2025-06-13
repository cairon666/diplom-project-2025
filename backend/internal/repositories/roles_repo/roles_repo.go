package roles_repo

import (
	"context"
	"errors"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/cairon666/vkr-backend/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RolesRepo struct {
	query    *dbqueries.Queries
	dbClient postgres.PostgresClient
}

func NewRolesRepo(query *dbqueries.Queries, dbClient postgres.PostgresClient) *RolesRepo {
	return &RolesRepo{
		query:    query,
		dbClient: dbClient,
	}
}

func (r *RolesRepo) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []int32) ([]models.Permission, error) {
	perms, err := r.query.GetPermissionsByRoleIDs(ctx, roleIDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, err
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

		return nil, err
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

// GetRolesByUserID получает список ролей пользователя по его UUID.
func (r *RolesRepo) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	roles, err := r.query.GetRolesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, err
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

// HasUserPermission проверяет, есть ли у пользователя разрешение с заданным именем.
func (r *RolesRepo) HasUserPermission(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	return r.query.HasUserPermission(ctx, dbqueries.HasUserPermissionParams{
		UserID: userID,
		Name:   permissionName,
	})
}

func (r *RolesRepo) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int32) error {
	err := r.query.AssignRoleToUser(ctx, dbqueries.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RolesRepo) GetRoleByName(ctx context.Context, name string) (models.Role, error) {
	role, err := r.query.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Role{}, apperrors.NotFound()
		}

		return models.Role{}, err
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

		return nil, err
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

func (r *RolesRepo) GetPermissionsByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Permission, error) {
	perms, err := r.query.GetPermissionsByExternalAppID(ctx, externalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound()
		}

		return nil, err
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

func (r *RolesRepo) AssignRoleToExternalApp(ctx context.Context, externalID uuid.UUID, roleID int32) error {
	err := r.query.AddRoleToExternalApp(ctx, dbqueries.AddRoleToExternalAppParams{
		ExternalAppID: externalID,
		RoleID:        roleID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RolesRepo) AssignRolesToExternalApp(ctx context.Context, externalID uuid.UUID, roleIDs []int32) error {
	tx, err := r.dbClient.Begin(ctx) // создаём транзакцию
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	for _, roleID := range roleIDs {
		err = r.query.AddRoleToExternalApp(ctx, dbqueries.AddRoleToExternalAppParams{
			ExternalAppID: externalID,
			RoleID:        roleID,
		})
		if err != nil {
			return err // defer откатит транзакцию
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit(ctx)

	return err
}
