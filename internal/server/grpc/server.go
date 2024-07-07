package grpc

import (
	"context"

	pb "github.com/artem-benda/monitor/internal/grpc/mon"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricsGrpsServer struct {
	pb.UnimplementedMonitorServiceServer
	Storage storage.Storage
	DBPool  *pgxpool.Pool
}

func (s *MetricsGrpsServer) GetMetric(c context.Context, req *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	switch {
	case req.Key.Id != "" && req.Key.Type == pb.MetricKey_COUNTER:
		cnt, ok, err := service.GetCounterMetric(c, s.Storage, req.Key.Id)
		if err != nil {
			return nil, status.Error(codes.Internal, "")
		}
		if !ok {
			return nil, status.Error(codes.NotFound, "")
		}
		return &pb.GetMetricResponse{Metric: &pb.MetricValue{MetricId: req.Key.Id, Value: &pb.MetricValue_Counter{Counter: cnt}}}, nil
	case req.Key.Id != "" && req.Key.Type == pb.MetricKey_GAUGE:
		val, ok, err := service.GetGaugeMetric(c, s.Storage, req.Key.Id)
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
	if req.Metric.MetricId == "" {
		return nil, status.Error(codes.NotFound, "")
	}

	switch req.Metric.Value.(type) {
	case *pb.MetricValue_Gauge:
		newGaugeVal, err := service.UpdateAndGetGaugeMetric(c, s.Storage, req.Metric.MetricId, req.Metric.GetGauge())
		if err != nil {
			return nil, status.Error(codes.Internal, "")
		}
		return &pb.UpdateMetricResponse{Metric: &pb.MetricValue{MetricId: req.Metric.MetricId, Value: &pb.MetricValue_Gauge{Gauge: newGaugeVal}}}, nil
	case *pb.MetricValue_Counter:
		newCounterVal, err := service.UpdateAndGetCounterMetric(c, s.Storage, req.Metric.MetricId, req.Metric.GetCounter())
		if err != nil {
			return nil, status.Error(codes.Internal, "")
		}
		return &pb.UpdateMetricResponse{Metric: &pb.MetricValue{MetricId: req.Metric.MetricId, Value: &pb.MetricValue_Counter{Counter: newCounterVal}}}, nil
	}

	return nil, status.Error(codes.Internal, "")
}

func (s *MetricsGrpsServer) UpdateMetricsBatch(c context.Context, req *pb.UpdateMetricsBatchRequest) (*emptypb.Empty, error) {
	logger.Log.Debug("MakeUpdateBatchJSONHandler, got metrics", zap.Int("count", len(req.Metrics)))

	if len(req.Metrics) == 0 {
		return &emptypb.Empty{}, nil
	}

	models := make([]model.MetricKeyWithValue, len(req.Metrics))

	for _, m := range req.Metrics {
		logger.Log.Debug("Adding metric...", zap.String("ID", m.MetricId))

		if m.MetricId == "" {
			return nil, status.Error(codes.NotFound, "")
		}

		switch m.Value.(type) {
		case *pb.MetricValue_Gauge:
			models = append(models, model.MetricKeyWithValue{Kind: model.GaugeKind, Name: m.MetricId, Gauge: m.GetGauge()})
			continue
		case *pb.MetricValue_Counter:
			models = append(models, model.MetricKeyWithValue{Kind: model.CounterKind, Name: m.MetricId, Counter: m.GetCounter()})
			continue
		}

		return nil, status.Error(codes.Internal, "")
	}

	err := service.UpdateMetrics(c, s.Storage, models)

	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &emptypb.Empty{}, nil
}

func (s *MetricsGrpsServer) PingDB(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	if s.DBPool == nil {
		return nil, status.Error(codes.Internal, "db connection is unused")
	}

	var ping string
	err := s.DBPool.QueryRow(context.Background(), "SELECT 'ping'").Scan(&ping)
	if err == nil && ping == "ping" {
		logger.Log.Debug("Executed ping command successfully")
		return &emptypb.Empty{}, nil
	} else {
		logger.Log.Debug("Executed ping command with error", zap.Error(err))
		return nil, status.Error(codes.Internal, "error executing db ping")
	}
}
