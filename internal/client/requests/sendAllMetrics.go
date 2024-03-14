package requests

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

func SendAllMetrics(resty *resty.Client, metrics map[model.MetricKey]model.MetricValue) error {
	dtos := make(dto.MetricsBatch, len(metrics))

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

	resp, err := resty.R().
		SetBody(b.Bytes()).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		Post("/updates/")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, status code = %d", resp.StatusCode())
	}
	return nil

}
