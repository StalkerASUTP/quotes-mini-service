package save

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"quotes-mini-service/internal/quote"
	"quotes-mini-service/pkg/res"
	"quotes-mini-service/pkg/sl"
	"strings"
)

type Request struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type QuoteSaver interface {
	Save(authorSave, quoteSave string) (*quote.Quote, error)
}

func New(log *slog.Logger, save QuoteSaver) http.HandlerFunc {
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
		if err = validate(req); err != nil {
			log.Error(err.Error(), sl.Err(err))
			res.Json(w, res.Error(err.Error()), http.StatusBadRequest)
			return
		}
		newQuote, err := save.Save(req.Author, req.Quote)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				log.Error("entry already exists", sl.Err(err))
				res.Json(w, res.Error("entry already exists"), http.StatusConflict)
				return
			}
			log.Error("failed to add quote", sl.Err(err))
			res.Json(w, res.Error("failed to add quote"), http.StatusInternalServerError)
			return
		}
		log.Info("quote added", slog.Int64("id", int64(newQuote.ID)))
		res.Json(w, newQuote, http.StatusCreated)
	}
}

func validate(req Request) error {
	switch {
	case req.Author == "" && req.Quote == "":
		return fmt.Errorf("author and quote are required")
	case req.Author == "":
		return fmt.Errorf("author is required")
	case req.Quote == "":
		return fmt.Errorf("quote is required")
	default:
		return nil
	}
}
