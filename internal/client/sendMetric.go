package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

func sendMetric(resty *resty.Client, apiURL string, kind string, name string, strVal string) error {
	rootURL, _ := strings.CutSuffix(apiURL, "/")
	url := fmt.Sprintf("%s/update/{metricKind}/{metricName}/{strVal}", rootURL)

	resp, err := resty.R().
		SetPathParams(map[string]string{
			"metricKind": kind,
			"metricName": name,
			"strVal":     strVal,
		}).
		SetBody("").
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, status code = %d", resp.StatusCode())
	}
	return nil
}
