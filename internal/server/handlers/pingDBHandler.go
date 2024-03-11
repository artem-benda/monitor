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

		var i int
		err := dbpool.QueryRow(context.Background(), "SELECT 1").Scan(&i)
		if i == 1 && err != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
