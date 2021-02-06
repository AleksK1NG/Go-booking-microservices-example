package main

import (
	"log"
	"os"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/config"
)

func main() {
	log.Println("Starting hotels microservice")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	log.Printf("CFG: %-v", cfg)
}
