package postgres

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

// Builder is a goqu postgres dialect builder.
var Builder = goqu.Dialect("postgres") //nolint:gochecknoglobals
