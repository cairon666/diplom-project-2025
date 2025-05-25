package role_service

import (
	"context"
	"fmt"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type RoleRepos interface {
	GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Permission, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int32) error
	GetRoleByName(ctx context.Context, name string) (models.Role, error)
	GetRolesByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Role, error)
	GetPermissionsByExternalAppID(ctx context.Context, externalID uuid.UUID) ([]models.Permission, error)
	AssignRoleToExternalApp(ctx context.Context, externalID uuid.UUID, roleID int32) error
	AssignRolesToExternalApp(ctx context.Context, externalID uuid.UUID, roleIDs []int32) error
}

type RoleService struct {
	rolesRepo RoleRepos
}

func NewRoleService(rolesRepo RoleRepos) *RoleService {
	return &RoleService{
		rolesRepo: rolesRepo,
	}
}

func (roleService *RoleService) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	perms, err := roleService.rolesRepo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p.Name == permission {
			return true, nil
		}
	}

	return false, nil
}

func (roleService *RoleService) HasOneOfPermissions(ctx context.Context, userID uuid.UUID, permissions []string) (bool, error) {
	perms, err := roleService.rolesRepo.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		for _, permission := range permissions {
			if p.Name == permission {
				return true, nil
			}
		}
	}

	return false, nil
}

func (roleService *RoleService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	role, err := roleService.rolesRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return err
	}

	if err := roleService.rolesRepo.AssignRoleToUser(ctx, userID, role.ID); err != nil {
		return err
	}

	return nil
}

func (roleService *RoleService) GetPermissionsByExternalAppUserID(ctx context.Context, userID uuid.UUID) ([]models.Permission, error) {
	return roleService.rolesRepo.GetPermissionsByExternalAppID(ctx, userID)
}

func (roleService *RoleService) GetRolesByExternalAppUserID(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	return roleService.rolesRepo.GetRolesByExternalAppID(ctx, userID)
}

func (roleService *RoleService) AssignRoleToExternalApp(ctx context.Context, externalID uuid.UUID, roleName string) error {
	role, err := roleService.rolesRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return err
	}

	if err := roleService.rolesRepo.AssignRoleToExternalApp(ctx, externalID, role.ID); err != nil {
		return err
	}

	return nil
}

func (roleService *RoleService) AssignRolesToExternalApp(ctx context.Context, externalID uuid.UUID, roleNames []string) error {
	roles := make([]int32, len(roleNames))
	for i, roleName := range roleNames {
		role, err := roleService.rolesRepo.GetRoleByName(ctx, roleName)
		if err != nil {
			return fmt.Errorf("failed to get role %s: %w", roleName, err)
		}

		roles[i] = role.ID
	}

	if err := roleService.rolesRepo.AssignRolesToExternalApp(ctx, externalID, roles); err != nil {
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}
