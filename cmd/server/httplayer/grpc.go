package httplayer

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"fmt"
	"log"

	app "github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
	pb "github.com/monsterr00/metric-service.gittest_client/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	appRepo app.App
}

func (s *MetricsServer) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse
	var metric models.Metric
	var err error

	metric.ID = in.Metric.Id
	metric.MType = in.Metric.MType
	metric.Delta = &in.Metric.Delta
	metric.Value = &in.Metric.Value

	switch metric.MType {
	case gaugeMetricType:
		err = s.appRepo.AddMetric(ctx, metric)
		if err != nil {
			return nil, status.Errorf(codes.Canceled, "Server: add metric error")
		}

	case counterMetricType:
		savedMetric, err := s.appRepo.Metric(ctx, metric.ID, metric.MType)
		if err == nil {
			var counter int64

			if savedMetric.Delta == nil {
				counter = 0
			} else {
				counter = *savedMetric.Delta
			}

			if metric.Delta != nil {
				counter += *metric.Delta
				metric.Delta = &counter
			}
		}
		err = s.appRepo.AddMetric(ctx, metric)
		if err != nil {
			return nil, status.Errorf(codes.Canceled, "Server: add metric error")
		}
	default:
		return nil, status.Errorf(codes.Canceled, "Wrong metric type")
	}

	log.Printf("Metric stored")

	return &response, nil
}

func (api *httpAPI) startGrpc() {
	pb.RegisterMetricsServer(api.grpc, &MetricsServer{appRepo: api.app})

	fmt.Println("Сервер gRPC начал работу")
}
