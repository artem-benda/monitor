package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
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

func (c *signWriter) WriteSigAndBody() {
	b := c.buf.Bytes()
	signature := Sign(b, c.key)
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	c.w.Header().Add(HashHeader, signatureBase64)
	c.w.Write(b)
}

func Sign(b []byte, signingKey []byte) []byte {
	h := hmac.New(sha256.New, signingKey)
	return h.Sum(b)
}

func Verify(b []byte, signature []byte, signingKey []byte) bool {
	h := hmac.New(sha256.New, signingKey)
	expectedSignature := h.Sum(b)
	return hmac.Equal(signature, expectedSignature)
}

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
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if !Verify(b, signature, signingKey) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
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
