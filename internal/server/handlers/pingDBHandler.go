package handlers

import (
	"context"
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func MakePingDatabaseHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ping string
		err := dbpool.QueryRow(context.Background(), "SELECT 'ping'").Scan(&ping)
		if err == nil && ping == "ping" {
			logger.Log.Debug("Executed ping command successfully")
			w.WriteHeader(http.StatusOK)
		} else {
			logger.Log.Debug("Executed ping command with error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
