package indentity

import (
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// AuthClaims Универсальная структура для авторизационных данных.
type AuthClaims struct {
	UserID      uuid.UUID
	Roles       []models.Role
	Permissions []models.Permission
}

// HasRole Проверяет наличие роли в списке ролей.
func (ac *AuthClaims) HasRole(roleName string) bool {
	for _, role := range ac.Roles {
		if role.Name == roleName {
			return true
		}
	}

	return false
}

// HasPermission Проверяет наличие разрешения в списке разрешений.
func (ac *AuthClaims) HasPermission(permissionsName ...string) bool {
	for _, permission := range ac.Permissions {
		for _, permissionName := range permissionsName {
			if permission.Name == permissionName {
				return true
			}
		}
	}

	return false
}
