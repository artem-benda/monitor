package handlers

import (
	"net/http"
	"strings"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
)

// MakeUpdatePathHandler - создать обработчик метода обновления значения указанной метрики из url path
func MakeUpdatePathHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain")

		switch metricKind, metricName, strVal := extractUpdatePathParams(r.URL.Path); {
		case model.ValidMetricKind(metricKind) && metricKind != "" && metricName != "":
			if err := service.UpdateMetric(r.Context(), store, metricKind, metricName, strVal); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		case r.Method != http.MethodPost:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case !model.ValidMetricKind(metricKind):
			w.WriteHeader(http.StatusBadRequest)
		case metricName == "":
			w.WriteHeader(http.StatusNotFound)
		case strVal == "":
			w.WriteHeader(http.StatusUnprocessableEntity)

		default:
			w.WriteHeader(http.StatusNotImplemented)
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
