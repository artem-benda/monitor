package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/service"
	"github.com/artem-benda/monitor/internal/storage"
)

func MakeUpdateHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("UpdateHandler, method = %s, path = %s", r.Method, r.URL.Path)
		w.Header().Add("Content-type", "text/plain")

		switch metricKind, metricName, strVal := extractUpdatePathParams(r.URL.Path); {
		case model.ValidMetricKind(metricKind) && metricKind != "" && metricName != "":
			if err := service.UpdateMetric(storage.Store, metricKind, metricName, strVal); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, "Bad metric value", http.StatusBadRequest)
			}
		case r.Method != http.MethodPost:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		case !model.ValidMetricKind(metricKind):
			http.Error(w, "Metric type not supported", http.StatusBadRequest)
		case metricName == "":
			http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
		case strVal == "":
			http.Error(w, "Invalid parameters values", http.StatusUnprocessableEntity)

		default:
			http.Error(w, "Method unimplemented", http.StatusNotImplemented)
		}
	}
}

func extractUpdatePathParams(urlPath string) (string, string, string) {
	var metricKind, metricName, strVal string

	params := strings.Split(strings.TrimPrefix(urlPath, "/update/"), "/")
	log.Printf("params: %s", params)

	switch {
	case len(params) >= 3:
		strVal = params[2]
		fallthrough
	case len(params) >= 2:
		metricName = params[1]
		fallthrough
	case len(params) >= 1:
		metricKind = params[0]
	}
	return metricKind, metricName, strVal
}
