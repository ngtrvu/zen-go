package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	Path string
	Port string
}

func NewMetricServer(path, port string) *MetricsServer {
	return &MetricsServer{Path: path, Port: port}
}

// TODO: using context to stop the metrics server
func (s *MetricsServer) Start(ctx context.Context) error {
	http.Handle(s.Path, promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%s", s.Port), nil)
}
