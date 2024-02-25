package requests

import (
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/go-resty/resty/v2"
)

func sendMetric(resty *resty.Client, dto dto.Metrics) error {
	resp, err := resty.R().
		SetBody(dto).
		Post("/update/")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, status code = %d", resp.StatusCode())
	}
	return nil
}
