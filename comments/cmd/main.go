package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AleksK1NG/hotels-mocroservices/comments/config"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/postgres"
)

func main() {
	log.Println("StartingComments service")

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

	appLogger.Infof("CFG: %-v", cfg)

	appLogger.Infof("%-v", pgxConn.Stat())

	http.ListenAndServe(":5555", nil)
}
