package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/config"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/pkg/logger"
)

// Server
type Server struct {
	logger    logger.Logger
	cfg       *config.Config
	redisConn *redis.Client
}

// NewServer
func NewSessionsServer(logger logger.Logger, cfg *config.Config, redisConn *redis.Client) *Server {
	return &Server{logger: logger, cfg: cfg, redisConn: redisConn}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	router := echo.New()
	router.GET("/api/v1/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		if err := router.Start(s.cfg.Metrics.URL); err != nil {
			s.logger.Errorf("router.Start metrics: %v", err)
			cancel()
		}
		s.logger.Infof("Metrics available on: %v", s.cfg.Metrics.URL)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	if err := router.Shutdown(ctx); err != nil {
		s.logger.Errorf("Metrics router.Shutdown: %v", err)
	}

	return nil
}
