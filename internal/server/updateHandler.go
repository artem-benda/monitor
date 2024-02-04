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
		switch params := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/"); {
		case r.Method != http.MethodPost:
			{
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case len(params) > 0 && !model.ValidMetricKind(params[0]):
			{
				http.Error(w, "Metric type not supported", http.StatusBadRequest)
			}
		case len(params) < 2 || params[1] == "":
			{
				http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
			}
		case len(params) < 3 || params[2] == "":
			{
				http.Error(w, "Invalid parameters values", http.StatusUnprocessableEntity)
			}
		case len(params) == 3 && model.ValidMetricKind(params[0]) && params[1] != "":
			{
				if err := service.UpdateMetric(storage.Store, params[0], params[1], params[2]); err == nil {
					w.WriteHeader(http.StatusOK)
				} else {
					http.Error(w, "Bad metric value", http.StatusBadRequest)
				}
			}
		default:
			{
				http.Error(w, "Method unimplemented", http.StatusNotImplemented)
			}
		}
	}
}
