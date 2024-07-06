package grpc

import (
	"context"
	"net/http"

	pb "github.com/artem-benda/monitor/internal/grpc/mon"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricsGrpsServer struct {
	pb.UnimplementedMonitorServiceServer
	storage storage.Storage
	dbpool  *pgxpool.Pool
}

func (s *MetricsGrpsServer) GetMetric(c context.Context, req *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	switch {
	case req.Key.Id != "" && req.Key.Type == pb.MetricKey_COUNTER:
		cnt, ok, err := service.GetCounterMetric(c, s.storage, req.Key.Id)
		if err != nil {
			return nil, status.Error(codes.Internal, "")
		}
		if !ok {
			return nil, status.Error(codes.NotFound, "")
		}
		return &pb.GetMetricResponse{Metric: &pb.MetricValue{MetricId: req.Key.Id, Value: &pb.MetricValue_Counter{Counter: cnt}}}, nil
	case req.Key.Id != "" && req.Key.Type == pb.MetricKey_GAUGE:
		val, ok, err := service.GetGaugeMetric(c, s.storage, req.Key.Id)
		if err != nil {
			return nil, status.Error(codes.Internal, "")
		}
		if !ok {
			return nil, status.Error(codes.NotFound, "")
		}
		return &pb.GetMetricResponse{Metric: &pb.MetricValue{MetricId: req.Key.Id, Value: &pb.MetricValue_Gauge{Gauge: val}}}, nil
	default:
		{
			return nil, status.Error(codes.NotFound, "")
		}
	}
}

func (s *MetricsGrpsServer) UpdateMetric(c context.Context, req *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	switch {
	case req.Metric.MetricKey.Id == "":
		return nil, status.Error(codes.NotFound, "")
	case req.Metric.MetricKey.Type == pb.MetricKey_GAUGE && (req.Metric.Value == nil):
		return nil, status.Error(codes.InvalidArgument, "")
	case metrics.MType == model.CounterKind && metrics.Delta == nil:
		w.WriteHeader(http.StatusBadRequest)
		return
	case metrics.MType == model.GaugeKind:
		var err error
		*metrics.Value, err = service.UpdateAndGetGaugeMetric(r.Context(), store, metrics.ID, *metrics.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _, err = easyjson.MarshalToHTTPResponseWriter(metrics, w)
			if err != nil {
				logger.Log.Error("Could not write json body", zap.Error(err))
			}
		}
	case metrics.MType == model.CounterKind:
		var err error
		*metrics.Delta, err = service.UpdateAndGetCounterMetric(r.Context(), store, metrics.ID, *metrics.Delta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _, err = easyjson.MarshalToHTTPResponseWriter(metrics, w)
			if err != nil {
				logger.Log.Error("Could not write json body", zap.Error(err))
			}
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "")
	}
}

func (s *MetricsGrpsServer) UpdateMetricsBatch(context.Context, *pb.UpdateMetricsBatchRequest) (*emptypb.Empty, error) {

}

func (s *MetricsGrpsServer) PingDB(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	if s.dbpool == nil {
		return nil, status.Error(codes.Internal, "db connection is unused")
	}

	var ping string
	err := s.dbpool.QueryRow(context.Background(), "SELECT 'ping'").Scan(&ping)
	if err == nil && ping == "ping" {
		logger.Log.Debug("Executed ping command successfully")
		return &emptypb.Empty{}, nil
	} else {
		logger.Log.Debug("Executed ping command with error", zap.Error(err))
		return nil, status.Error(codes.Internal, "error executing db ping")
	}
}

func mapToMetricType(t pb.MetricKey_MetricType) string {
	switch t {
	case pb.MetricKey_COUNTER:
		return model.CounterKind
	case pb.MetricKey_GAUGE:
		return model.GaugeKind
	default:
		return ""
	}
}
