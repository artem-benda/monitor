package handlers

import (
	"net/http"
	"strings"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func MakeUpdatePathHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain")

		switch metricKind, metricName, strVal := extractUpdatePathParams(r.URL.Path); {
		case model.ValidMetricKind(metricKind) && metricKind != "" && metricName != "":
			if err := service.UpdateMetric(store, metricKind, metricName, strVal); err == nil {
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
