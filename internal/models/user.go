package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                     uuid.UUID
	Email                  string
	FirstName              string
	LastName               string
	IsRegistrationComplete bool
	CreatedAt              time.Time
}

func NewUser(ID uuid.UUID, email string, firstName string, lastName string, isRegistrationComplete bool, createdAt time.Time) *User {
	return &User{
		ID:                     ID,
		Email:                  email,
		FirstName:              firstName,
		LastName:               lastName,
		IsRegistrationComplete: isRegistrationComplete,
		CreatedAt:              createdAt,
	}
}

type UserPassword struct {
	UserId       uuid.UUID
	Salt         string
	PasswordHash string
}

func NewUserPassword(userId uuid.UUID, salt string, passwordHash string) UserPassword {
	return UserPassword{
		UserId:       userId,
		Salt:         salt,
		PasswordHash: passwordHash,
	}
}

const TelegramProviderName = "telegram"

type AuthProvider struct {
	ID             uuid.UUID
	UserId         uuid.UUID
	ProviderName   string
	ProviderUserId int64
	CreatedAt      time.Time
}

func NewAuthProvider(ID uuid.UUID, userId uuid.UUID, providerName string, providerUserId int64, createdAt time.Time) AuthProvider {
	return AuthProvider{
		ID:             ID,
		UserId:         userId,
		ProviderName:   providerName,
		ProviderUserId: providerUserId,
		CreatedAt:      createdAt,
	}
}
