package requests

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func sendMetric(resty *resty.Client, kind string, name string, strVal string) error {
	resp, err := resty.R().
		SetPathParams(map[string]string{
			"metricKind": kind,
			"metricName": name,
			"strVal":     strVal,
		}).
		SetBody("").
		Post("/update/{metricKind}/{metricName}/{strVal}")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, status code = %d", resp.StatusCode())
	}
	return nil
}
