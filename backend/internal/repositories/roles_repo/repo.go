package roles_repo

import (
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/cairon666/vkr-backend/pkg/postgres"
)

type RolesRepo struct {
	query    *dbqueries.Queries
	dbClient postgres.PostgresClient
}

func NewRolesRepo(dbClient postgres.PostgresClient) *RolesRepo {
	return &RolesRepo{
		query:    dbqueries.New(dbClient),
		dbClient: dbClient,
	}
}
