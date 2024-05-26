package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"go.uber.org/zap"
)

const HashHeader = "HashSHA256"

type signWriter struct {
	w   http.ResponseWriter
	key []byte
	buf *bytes.Buffer
}

func newSignWriter(w http.ResponseWriter, signingKey []byte) *signWriter {
	return &signWriter{
		w:   w,
		key: signingKey,
		buf: &bytes.Buffer{},
	}
}

func (c *signWriter) Header() http.Header {
	return c.w.Header()
}

// Пишем в буффер, чтобы иметь возможность записать заголовок на основе тела запроса
func (c *signWriter) Write(p []byte) (int, error) {
	return c.buf.Write(p)
}

func (c *signWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *signWriter) WriteSigAndBody() error {
	b := c.buf.Bytes()
	signature, err := Sign(b, c.key)
	if err != nil {
		logger.Log.Debug("Error signing response", zap.Error(err))
		return err
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	c.w.Header().Add(HashHeader, signatureBase64)
	_, err = c.w.Write(b)
	return err
}

// Sign - Вычислить подпись для слайса байт с указанным ключом с использованием алгоритма hmac-sha256
func Sign(b []byte, signingKey []byte) ([]byte, error) {
	h := hmac.New(sha256.New, signingKey)
	_, err := h.Write(b)
	if err != nil {
		logger.Log.Debug("Error creating signature", zap.Error(err))
	}
	return h.Sum(nil), nil
}

// Verify - Проверить подпись для слайса байт с использованием алгоритма hmac-sha256
func Verify(b []byte, signature []byte, signingKey []byte) (bool, error) {
	expectedSignature, err := Sign(b, signingKey)
	if err != nil {
		logger.Log.Debug("Error verifying signature", zap.Error(err))
	}
	logger.Log.Debug("Signatures", zap.String("actual", string(signature)), zap.String("expected", string(expectedSignature)))
	return hmac.Equal(signature, expectedSignature), nil
}

// CreateVerifyAndSignMiddleware - создать новый middleware для проверки и формирования подписи с указанным ключом
func CreateVerifyAndSignMiddleware(signingKey []byte) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
			// который будем передавать следующей функции
			ow := w

			// проверяем, что клиент прислал запрос с подписью
			signatureHeaderValue := r.Header.Get(HashHeader)
			if len(signatureHeaderValue) > 0 && len(signingKey) > 0 {
				signature, err := base64.StdEncoding.DecodeString(signatureHeaderValue)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// меняем тело запроса на новое
				b, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Log.Debug("Error reading body", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if ok, err := Verify(b, signature, signingKey); err != nil {
					logger.Log.Debug("Error verifying signature", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				} else if !ok {
					logger.Log.Debug("Invalid signature")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(b))
			}

			if len(signingKey) > 0 {
				sw := newSignWriter(w, signingKey)
				ow = sw
				// Вычисляем подпись, устанавливаем заголовок и записываем тело запроса из буффера
				defer sw.WriteSigAndBody()
			}

			// передаём управление хендлеру
			h.ServeHTTP(ow, r)
		})
	}
}
