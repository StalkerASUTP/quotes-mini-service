package main

import (
	"log/slog"
	"os"
	"quotes-mini-service/internal/config"
	"quotes-mini-service/internal/storage"
)

func main() {
	conf := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("initializing server", slog.String("address", conf.Address))
	log.Debug("logger debug mode enabled")
	db, err := storage.NewStorage(conf.Database)
	if err != nil {
		log.Error("failed to initialize storage", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}

}
