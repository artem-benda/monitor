package requests

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/go-resty/resty/v2"
)

func sendMetric(resty *resty.Client, dto dto.Metrics) error {
	var (
		json []byte
		err  error
	)

	if json, err = resty.JSONMarshal(dto); err != nil {
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
		Post("/update/")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, status code = %d", resp.StatusCode())
	}
	return nil
}
