package main

import (
	"log"
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/jaeger"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/postgres"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/redis"
)

func main() {
	log.Println("Starting user server")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s",
		cfg.GRPCServer.AppVersion,
		cfg.Logger.Level,
		cfg.GRPCServer.Mode,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.GRPCServer.AppVersion)

	pgxConn, err := postgres.NewPgxConn(cfg)
	if err != nil {
		appLogger.Fatal("cannot connect to postgres", err)
	}
	defer pgxConn.Close()

	redisClient := redis.NewRedisClient(cfg)
	appLogger.Info("Redis connected")

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	log.Printf("%-v", pgxConn.Stat())
	log.Printf("%-v", redisClient.PoolStats())

	http.ListenAndServe(":5001", nil)
}
