package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MakePingDatabaseHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dbpool == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ping string
		err := dbpool.QueryRow(context.Background(), "SELECT 'ping'").Scan(&ping)
		if err != nil && ping == "ping" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
