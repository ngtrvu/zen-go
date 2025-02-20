package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/ngtrvu/zen-go/log"
	"github.com/ngtrvu/zen-go/metrics"

	http_metrics "github.com/ngtrvu/zen-go/metrics/http"
)

type APIServer struct {
	HTTPServer http.Server
}

func NewRouter(AppName string) *chi.Mux {
	router := chi.NewRouter()

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// sentry middleware
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
	router.Use(sentryMiddleware.Handle)

	// add some base middlewares
	httpMetricsObserver := http_metrics.NewHTTPMetrics(AppName)
	router.Use(http_metrics.InboundMetricsMiddleware(httpMetricsObserver))
	router.Use(middleware.RealIP)

	router.Get("/healthy", GetHealthCheck)

	return router
}

func NewAPIServer(router *chi.Mux, httpServerConfig *HTTPServerConfig) *APIServer {
	apiServer := &APIServer{
		HTTPServer: http.Server{
			Addr:    fmt.Sprintf(":%d", httpServerConfig.Port),
			Handler: router,
		},
	}

	return apiServer
}

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, "OK")
}

func (apiServer *APIServer) Start(ctx context.Context) error {
	go func() {
		log.Info("starting metrics server at port 9090...")

		metricsServer := metrics.NewMetricServer("/metrics", "9090")
		err := metricsServer.Start(ctx)
		if err != nil {
			log.Info("metrics server return error: %v\n", err)
		}
	}()

	go func() {
		defer sentry.Flush(2 * time.Second)

		err := apiServer.HTTPServer.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("http server error: %v", err)
			} else {
				log.Info("stopped serving new connections")
			}
		}
	}()

	log.Info("listening and serving...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("received OS signal. Shutting down server...")
	apiServer.HTTPServer.Shutdown(ctx)
	log.Info("server stopped safely")
	return nil
}
