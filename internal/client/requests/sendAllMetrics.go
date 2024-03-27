package requests

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/client/errors"
	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/retry"
	"github.com/artem-benda/monitor/internal/signer"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

func SendAllMetrics(c *resty.Client, withRetry retry.RetryController, metrics map[model.MetricKey]model.MetricValue, signingKey []byte) error {
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
	if _, err := w.Write(json); err != nil {
		return err
	}
	w.Close()

	var resp *resty.Response

	err = withRetry.Run(func() (err error) {
		resp, err = sendBytes(c, b.Bytes(), signingKey)
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

func sendBytes(resty *resty.Client, b []byte, signingKey []byte) (*resty.Response, error) {
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
