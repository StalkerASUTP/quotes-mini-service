package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"quotes-mini-service/internal/quote"
	"quotes-mini-service/pkg/res"
	"quotes-mini-service/pkg/sl"
)

type Request struct {
	Author string `json:"author" validate:"required"`
	Quote  string `json:"quote" validate:"required"`
}
type QuoteSaver interface {
	Save(authorSave, quoteSave string) (*quote.Quote, error)
}

func New(log *slog.Logger, saver QuoteSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"
		log = log.With(
			slog.String("op", op),
		)
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			res.Json(w, res.Error("failed to decode request"), http.StatusBadRequest)
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		newQuote, err := saver.Save(req.Author, req.Quote)
		if err != nil {
			log.Error("failed to add quote", sl.Err(err))
			res.Json(w, res.Error("failed to add quote"), http.StatusInternalServerError)
			return
		}
		log.Info("quote added", slog.Int64("id", int64(newQuote.ID)))
		res.Json(w, newQuote, http.StatusCreated)
	}
}
