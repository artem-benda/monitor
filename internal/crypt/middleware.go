package crypt

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"
	"strings"

	"github.com/artem-benda/monitor/internal/logger"
	"go.uber.org/zap"
)

// NewDecryptMiddleware - создать middleware расшифрования запросов с заданным приватным ключом
func NewDecryptMiddleware(privateKey *rsa.PrivateKey) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// проверяем, что клиент отправил серверу зашифрованные данные
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsEncrypted := strings.Contains(contentEncoding, "encrypted")
			if sendsEncrypted {
				encBody, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Log.Debug("Error reading body", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// Рашифровываем тело запроса
				b, err := DecryptWithPrivateKey(encBody, privateKey)
				if err != nil {
					logger.Log.Debug("Error decrypting body", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				// меняем тело запроса на новое
				r.Body = io.NopCloser(bytes.NewBuffer(b))
			}

			// передаём управление хендлеру
			h.ServeHTTP(w, r)
		})
	}
}
