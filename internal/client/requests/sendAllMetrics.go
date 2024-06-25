package requests

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/client/errors"
	"github.com/artem-benda/monitor/internal/crypt"
	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/retry"
	"github.com/artem-benda/monitor/internal/signer"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

// SendAllMetrics - Отправить значения метрик на сервер
func SendAllMetrics(c *resty.Client, withRetry retry.RetryController, metrics map[model.MetricKey]model.MetricValue, signingKey []byte, rsaPublicKey *rsa.PublicKey) error {
	dtos := make(dto.MetricsBatch, 0)

	for k, v := range metrics {
		dtos = append(dtos, model.AsDto(k, v))
	}

	var (
		json []byte
		err  error
	)

	if json, err = easyjson.Marshal(dtos); err != nil {
		return err
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err = w.Write(json); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	var body []byte
	var isEncryptionEnabled bool
	if rsaPublicKey == nil {
		body = b.Bytes()
	} else {
		bytes, err := crypt.EncryptWithPublicKey(b.Bytes(), rsaPublicKey)
		if err != nil {
			logger.Log.Error("Unable to encrypt body", zap.Int("body_bytes", len(b.Bytes())), zap.Error(err))
			return err
		}
		body = bytes
		isEncryptionEnabled = true
	}

	var resp *resty.Response

	err = withRetry.Run(func() (err error) {
		resp, err = sendBytes(c, body, signingKey, isEncryptionEnabled)
		return
	})

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, status code = %d", resp.StatusCode())
	}
	return nil

}

func sendBytes(resty *resty.Client, b []byte, signingKey []byte, isEncrypted bool) (*resty.Response, error) {
	var signatureBase64 string
	if len(signingKey) > 0 {
		signature, err := signer.Sign(b, signingKey)
		if err != nil {
			logger.Log.Debug("Error signing metrics", zap.Error(err))
			return nil, err
		}
		signatureBase64 = base64.StdEncoding.EncodeToString(signature)
	}

	request := resty.R().
		SetBody(b).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	if signatureBase64 != "" {
		request = request.
			SetHeader("HashSHA256", signatureBase64)
	}
	if isEncrypted {
		request = request.
			SetHeader("Content-Encryption", "encrypted")
	}

	resp, err := request.
		Post("/updates/")

	if err != nil {
		logger.Log.Debug("Error sending metrics", zap.Error(err))
		return nil, errors.ErrNetwork{Err: err}
	}

	if resp.StatusCode() >= 500 {
		return nil, errors.ErrServerTemporary{Err: err}
	}

	return resp, err
}
