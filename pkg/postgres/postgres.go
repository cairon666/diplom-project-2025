package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

func NewPostgres(ctx context.Context, connString string) (PostgresClient, error) {
	pgx, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return pgx, nil
}
