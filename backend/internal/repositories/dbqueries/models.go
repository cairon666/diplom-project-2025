// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package dbqueries

import (
	"time"

	"github.com/google/uuid"
)

type AUTHPROVIDER struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	ProviderName   string
	CreatedAt      time.Time
	ProviderUserID int64
}

type DEVICE struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	DeviceName string
	CreatedAt  time.Time
}

type EXTERNALAPP struct {
	ID          uuid.UUID
	Name        string
	OwnerUserID uuid.UUID
	ApiKeyHash  string
	CreatedAt   time.Time
}

type EXTERNALAPPSROLE struct {
	ExternalAppID uuid.UUID
	RoleID        int32
}

type PERMISSION struct {
	ID   int32
	Name string
}

type ROLE struct {
	ID   int32
	Name string
}

type ROLEPERMISSION struct {
	RoleID       int32
	PermissionID int32
}

type USER struct {
	ID                     uuid.UUID
	Email                  *string
	FirstName              string
	LastName               string
	IsRegistrationComplete bool
	CreatedAt              time.Time
}

type USERPASSWORD struct {
	UserID       uuid.UUID
	PasswordHash string
	Salt         string
}

type USERROLE struct {
	UserID uuid.UUID
	RoleID int32
}
