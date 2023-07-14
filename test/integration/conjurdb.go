// Package main provides integration tests for conjur service broker
package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type conjurdb struct {
	cfg config
}

func (*conjurdb) conjurResourceExists(id string) error {
	connStr := "postgres://postgres@postgres:5432/postgres"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer conn.Close(context.Background())

	q, err := conn.Exec(context.Background(), "SELECT 1 FROM resources WHERE resource_id=$1", id)
	if err != nil {
		return fmt.Errorf("unable to query database: %w", err)
	}
	if q.RowsAffected() != 1 {
		return fmt.Errorf("resource with id: %s not found in the database", id)
	}
	return nil
}
