package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server
type Server struct {
	echo      *echo.Echo
	logger    logger.Logger
	cfg       *config.Config
	redisConn *redis.Client
	pgxPool   *pgxpool.Pool
}

// NewServer
func NewServer(logger logger.Logger, cfg *config.Config, redisConn *redis.Client, pgxPool *pgxpool.Pool) *Server {
	return &Server{logger: logger, cfg: cfg, redisConn: redisConn, pgxPool: pgxPool, echo: echo.New()}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	s.logger.Info("Server Exited Properly")

	if err := s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	return nil
}
