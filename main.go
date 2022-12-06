package main

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./externalscaler/externalscaler.proto

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/imroc/req/v3"

	pb "github.com/LCOGT/keda-celery-flower-scaler/externalscaler"
)

var (
	address = flag.String("address", ":0", "host:port to listen on")
)

const (
	RunningAndPendingTasksMetricName = "runningAndPendingTasks"
)

func main() {
	flag.Parse()

	grpcServer := grpc.NewServer()

	s := &CeleryScaler{
		client: req.C(),
	}

	pb.RegisterExternalScalerServer(grpcServer, s)

	reflection.Register(grpcServer)

	l, err := net.Listen("tcp4", *address)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening on %v.", l.Addr())

	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

}

type CeleryScaler struct {
	pb.UnimplementedExternalScalerServer

	client *req.Client
}

func (self *CeleryScaler) IsActive(ctx context.Context, so *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	config, err := parseConfig(so)

	if err != nil {
		return nil, err
	}

	count, err := self.queryRunningAndPendingTasks(ctx, config.FlowerURL)

	if err != nil {
		return nil, err
	}

	resp := &pb.IsActiveResponse{
		Result: count >= config.ActivationThreshold,
	}

	return resp, nil
}

func (self *CeleryScaler) StreamIsActive(so *pb.ScaledObjectRef, stream pb.ExternalScaler_StreamIsActiveServer) error {
	return status.Error(codes.Unimplemented, "TODO")
}

func (self *CeleryScaler) GetMetricSpec(ctx context.Context, so *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	config, err := parseConfig(so)

	if err != nil {
		return nil, err
	}

	resp := &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{
			{
				MetricName: RunningAndPendingTasksMetricName,
				TargetSize: int64(config.DesiredMetricValue),
			},
		},
	}

	return resp, nil
}

func (self *CeleryScaler) GetMetrics(ctx context.Context, mr *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	config, err := parseConfig(mr.ScaledObjectRef)

	if err != nil {
		return nil, err
	}

	count, err := self.queryRunningAndPendingTasks(ctx, config.FlowerURL)

	if err != nil {
		return nil, err
	}

	resp := &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{
			{
				MetricName:  RunningAndPendingTasksMetricName,
				MetricValue: int64(count),
			},
		},
	}
	return resp, nil
}
