package handlers

import (
	"context"
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// MakePingDatabaseHandler - создать обработчик метода проверки статуса соединения с хранилищем
func MakePingDatabaseHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dbpool == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
