package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	traceutils "github.com/opentracing-contrib/go-grpc"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	grpcImg "github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/delivery/grpc"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/delivery/rabbitmq"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/publisher"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/repository"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/usecase"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
	imageService "github.com/AleksK1NG/hotels-mocroservices/images-microservice/proto/image"
)

type Server struct {
	logger logger.Logger
	cfg    *config.Config
	tracer opentracing.Tracer
}

func NewServer(logger logger.Logger, cfg *config.Config, tracer opentracing.Tracer) *Server {
	return &Server{logger: logger, cfg: cfg, tracer: tracer}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	imagePublisher, err := publisher.NewImagePublisher(s.cfg, s.logger)
	if err != nil {
		return errors.Wrap(err, "NewImagePublisher")
	}
	uploadedChan, err := imagePublisher.CreateExchangeAndQueue("images", "uploaded", "uploaded")
	if err != nil {
		return errors.Wrap(err, "imagePublisher.CreateExchangeAndQueue")
	}
	defer uploadedChan.Close()

	// Init repos, usecases, middlewares, interceptors
	imagePGRepo := repository.NewImagePGRepository()
	imageAWSRepo := repository.NewImageAWSRepository()
	imageUC := usecase.NewImageUseCase(imagePGRepo, imageAWSRepo, s.logger, imagePublisher)

	// Init consumers, publishers, grpc server, metrics
	imageConsumer := rabbitmq.NewImageConsumer(s.logger, s.cfg, imageUC)
	if err := imageConsumer.Dial(); err != nil {
		return errors.Wrap(err, " imageConsumer.Dial")
	}

	resizeChan, err := imageConsumer.CreateExchangeAndQueue("images", "resize", "images")
	if err != nil {
		return errors.Wrap(err, "CreateExchangeAndQueue")
	}
	defer resizeChan.Close()

	go func() {
		if err := imageConsumer.StartResizeConsumer(
			ctx,
			24,
			"resize",
			"resizeConsumer",
		); err != nil {
			s.logger.Errorf("StartResizeConsumer: %v", err)
			cancel()
		}
	}()

	l, err := net.Listen("tcp", s.cfg.GRPCServer.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.GRPCServer.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.GRPCServer.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.GRPCServer.MaxConnectionAge * time.Minute,
		Time:              s.cfg.GRPCServer.Timeout * time.Minute,
	}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			traceutils.OpenTracingServerInterceptor(s.tracer),
			// im.Logger,
		),
	)

	imgService := grpcImg.NewImageService(s.cfg, s.logger)
	imageService.RegisterImageServiceServer(server, imgService)
	grpc_prometheus.Register(server)

	go func() {
		s.logger.Infof("GRPC Server is listening on port: %v", s.cfg.GRPCServer.Port)
		s.logger.Fatal(server.Serve(l))
	}()

	if s.cfg.GRPCServer.Mode != "Production" {
		reflection.Register(server)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	server.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}