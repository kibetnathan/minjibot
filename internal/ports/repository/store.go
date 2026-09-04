// Package repository defines repository interfaces for all domain entities
// and provides SQL implementations backed by sqlc-generated queries.
package repository

import (
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
)

// SQLStore Def --
// Shared by all repositories
type SQLStore struct {
	queries postgres.Querier
}

func NewSQLStore(queries postgres.Querier) *SQLStore {
	return &SQLStore{queries: queries}
}
