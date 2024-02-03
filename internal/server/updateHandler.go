package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/service"
	"github.com/artem-benda/monitor/internal/storage"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateHandler, method = %s, path = %s", r.Method, r.URL.Path)
	switch params := strings.Split(r.URL.Path, "/"); {
	case r.Method != http.MethodPost:
		{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case !strings.HasPrefix(r.Header.Get("Content-type"), "text/plain"):
		{
			http.Error(w, "Content type not supported", http.StatusBadRequest)
		}
	case len(params) == 3 && !model.ValidMetricKind(params[0]):
		{
			http.Error(w, "Metric type not supported", http.StatusBadRequest)
		}
	case len(params) == 3 && params[1] == "":
		{
			http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
		}
	case len(params) != 3:
		{
			http.Error(w, "Invalid parameters values", http.StatusUnprocessableEntity)
		}
	case len(params) == 3 && model.ValidMetricKind(params[0]) && params[1] != "":
		{
			if err := service.UpdateMetric(storage.Store, params[0], params[1], params[2]); err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, "Bad metric value", http.StatusUnprocessableEntity)
			}
		}
	default:
		{
			http.Error(w, "Method unimplemented", http.StatusNotImplemented)
		}
	}
}
