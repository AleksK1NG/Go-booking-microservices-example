package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/config"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
)

const (
	certFile          = "ssl/server.crt"
	keyFile           = "ssl/server.pem"
	maxHeaderBytes    = 1 << 20
	userCachePrefix   = "users:"
	userCacheDuration = time.Minute * 15
)

// Server
type server struct {
	echo      *echo.Echo
	logger    logger.Logger
	cfg       *config.Config
	redisConn *redis.Client
	tracer    opentracing.Tracer
}

// NewServer
func NewServer(logger logger.Logger, cfg *config.Config, redisConn *redis.Client, tracer opentracing.Tracer) *server {
	return &server{echo: echo.New(), logger: logger, cfg: cfg, redisConn: redisConn, tracer: tracer}
}

func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// validate := validator.New()
	v1 := s.echo.Group("/api/v1")
	hotels := v1.Group("/hotels")

	hotels.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Ok")
	})
	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	go func() {
		router := echo.New()
		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		s.logger.Infof("Metrics server is running on port: %s", s.cfg.Metrics.Port)
		if err := router.Start(s.cfg.Metrics.Port); err != nil {
			s.logger.Error(err)
			cancel()
		}
	}()

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.HttpServer.Port)
		s.echo.Server.ReadTimeout = time.Second * s.cfg.HttpServer.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.HttpServer.WriteTimeout
		s.echo.Server.MaxHeaderBytes = maxHeaderBytes
		if err := s.echo.StartTLS(s.cfg.HttpServer.Port, certFile, keyFile); err != nil {
			s.logger.Fatalf("Error starting TLS Server: ", err)
		}
	}()

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.HttpServer.PprofPort)
		if err := http.ListenAndServe(s.cfg.HttpServer.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	if err := s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "Server.Shutdown")
	}

	s.logger.Info("Server Exited Properly")

	return nil
}
