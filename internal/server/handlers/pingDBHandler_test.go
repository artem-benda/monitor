package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
)

func ExampleMakePingDatabaseHandler() {
	databaseDSN := os.Getenv("DATABASE_DSN")
	dbpool, err := pgxpool.New(context.Background(), databaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	handler := MakePingDatabaseHandler(dbpool)

	r := chi.NewRouter()
	r.Get("/", handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
