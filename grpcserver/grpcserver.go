package grpcserver

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/ngtrvu/zen-go/log"
	"github.com/ngtrvu/zen-go/metrics"
	grpc_metrics "github.com/ngtrvu/zen-go/metrics/grpc"
)

var (
	port = flag.Int("port", 8089, "The server port")
)

type StagGrpcServer struct {
	GrpcServer  *grpc.Server
	namespace   string
	serviceName string
}

func NewStagGrpcServer(ctx context.Context, ctxCancel context.CancelFunc, namespace string, serviceName string) *StagGrpcServer {
	server := new(StagGrpcServer)
	server.namespace = namespace
	server.serviceName = serviceName

	err := server.initialize(ctx, ctxCancel)
	if err != nil {
		log.Info("server return error: %v", err)
	}

	return server
}

func (s *StagGrpcServer) initialize(ctx context.Context, ctxCancel context.CancelFunc) error {
	// waiting for os signal
	go func() {
		defer ctxCancel()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		log.Info("received OS signal %v", <-sigChan)
	}()

	go func() {
		log.Info("starting app metrics server at port 9091...")
		metricsServer := metrics.NewMetricServer("/metrics", "9091")
		err := metricsServer.Start(ctx)
		if err != nil {
			log.Info("metrics server return error: %v", err)
		}
	}()

	metricsObserver := grpc_metrics.NewMetrics(s.namespace, s.serviceName)
	unaryInterceptor := grpc.ChainUnaryInterceptor(
		grpc_metrics.NewMetricsUnaryInterceptor(metricsObserver),
	)
	s.GrpcServer = grpc.NewServer(
		unaryInterceptor,
	)
	grpc_health_v1.RegisterHealthServer(s.GrpcServer, health.NewServer())

	return nil
}

func (s *StagGrpcServer) Serve(ctx context.Context, ctxCancel context.CancelFunc) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
		return err
	}

	go func() {
		defer ctxCancel()

		log.Info("running NewGrpcServer, port: %d...", *port)
		if err := s.GrpcServer.Serve(listener); err != nil {
			log.Fatal("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	s.GrpcServer.GracefulStop()
	log.Info("GRPC server stopped safely")

	return nil
}
