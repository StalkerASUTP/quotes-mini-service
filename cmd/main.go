package main

import (
	"log/slog"
	"net/http"
	"os"
	"quotes-mini-service/internal/config"
	"quotes-mini-service/internal/quote"
	deleting "quotes-mini-service/internal/quote/handlers/delete"
	"quotes-mini-service/internal/quote/handlers/get"
	"quotes-mini-service/internal/quote/handlers/save"
	"quotes-mini-service/internal/storage"
	"quotes-mini-service/pkg/sl"
)

func main() {
	conf := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("initializing server", slog.String("address", conf.Address))
	log.Debug("logger debug mode enabled")
	db, err := storage.NewStorage(conf.Database)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}
	queryRepository := quote.NewQuotesRepository(db)
	log.Info("initializing query repository")
	router := http.NewServeMux()
	router.HandleFunc("POST /quotes", save.New(log, queryRepository))
	router.HandleFunc("DELETE /quotes/{id}", deleting.New(log, queryRepository))
	router.HandleFunc("GET /quotes", get.AllParam(log, queryRepository))
	router.HandleFunc("GET /quotes/random", get.Random(log, queryRepository))

	log.Info("starting server", slog.String("address", conf.Address))
	server := http.Server{
		Addr:         conf.Address,
		Handler:      router,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
		IdleTimeout:  conf.IdleTimeout,
	}
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

}
