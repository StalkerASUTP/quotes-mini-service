package get

import (
	"log/slog"
	"net/http"
	"quotes-mini-service/internal/quote"
	"quotes-mini-service/pkg/res"
	"quotes-mini-service/pkg/sl"
)

type QuoteGetter interface {
	GetAllParam(author string) ([]quote.Quote, int, error)
	GetRandom() (*quote.Quote, error)
}
type GetWithParamResponse struct {
	Quotes []quote.Quote `json:"quotes"`
	Count  int           `json:"count" example:"100"`
}

func NewResponseWithParam(quotes []quote.Quote, count int) *GetWithParamResponse {
	return &GetWithParamResponse{
		Quotes: quotes,
		Count:  count,
	}
}
func AllParam(log *slog.Logger, get QuoteGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.All"
		log = log.With(
			slog.String("op", op),
		)
		author := r.URL.Query().Get("author")
		var (
			quotes []quote.Quote
			count  int
			err    error
		)
		quotes, count, err = get.GetAllParam(author)
		if err != nil {
			log.Error("internal server error", sl.Err(err))
			res.Json(w, res.Error("internal server error"), http.StatusInternalServerError)
			return
		}
		log.Info("slice of quotes getted")
		response := NewResponseWithParam(quotes, count)
		log.Info("response id ready", slog.Any("quotes", response), slog.Int("count", count))
		res.Json(w, response, http.StatusOK)
	}
}

func Random(log *slog.Logger, get QuoteGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.Random"
		log = log.With(
			slog.String("op", op),
		)
		randomQuote, err := get.GetRandom()
		if err != nil {
			log.Error("internal server error", sl.Err(err))
			res.Json(w, res.Error("internal server error"), http.StatusInternalServerError)
			return
		}
		log.Info("quote getted", slog.Any("res", randomQuote))
		res.Json(w, randomQuote, http.StatusOK)
	}
}
