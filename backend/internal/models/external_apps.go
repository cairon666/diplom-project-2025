package models

import (
	"time"

	"github.com/google/uuid"
)

type ExternalApp struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	Name       string
	APIKeyHash string
	CreatedAt  time.Time
}

func NewExternalApp(id, ownerID uuid.UUID, name string, apiKeyHash string, createdAt time.Time) ExternalApp {
	return ExternalApp{
		ID:         id,
		OwnerID:    ownerID,
		Name:       name,
		APIKeyHash: apiKeyHash,
		CreatedAt:  createdAt,
	}
}
