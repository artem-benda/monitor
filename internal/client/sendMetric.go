package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func SendMetric(httpClient *http.Client, apiURL string, kind string, name string, strVal string) error {
	rootURL, _ := strings.CutSuffix(apiURL, "/")
	url := fmt.Sprintf("%s/%s/%s/%s", rootURL, url.PathEscape(kind), url.PathEscape(name), url.PathEscape(strVal))

	resp, err := httpClient.Post(url, "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
