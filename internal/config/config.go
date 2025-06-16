package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database string
	HTTPServer
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var cfg Config
	addr := os.Getenv("APP_SERVER_ADDR")
	if addr == "" {
		log.Fatal("Empty server address")
	}
	cfg.Address = addr
	cfg.Timeout = parseDuration(os.Getenv("APP_TIMEOUT"))
	cfg.IdleTimeout = parseDuration(os.Getenv("APP_IDLE_TIMEOUT"))
	db := os.Getenv("APP_DATABASE")
	if db == "" {
		log.Fatal("Empty databases configration")
	}
	cfg.Database = db
	return &cfg
}

func parseDuration(times string) time.Duration {
	parsedTime, err := time.ParseDuration(times)
	if err != nil {
		parsedTime = time.Second * 10
	}
	return parsedTime
}
