package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newConnectionPool(databaseDSN string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), databaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}
