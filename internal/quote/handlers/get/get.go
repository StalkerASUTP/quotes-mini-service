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
		respose := GetWithParamResponse{
			Quotes: quotes,
			Count:  count,
		}
		log.Info("response id ready", slog.Any("quotes", respose), slog.Int("count", count))
		res.Json(w, respose, http.StatusOK)
	}
}

// func Param(log *slog.Logger, get QuoteGetter) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "handlers.get.Param"
// 		log = log.With(
// 			slog.String("op", op),
// 		)

// 		author := r.URL.Query().Get("author")
// 		if author == "" {
// 			log.Error("parameter author is empty")
// 			res.Json(w, res.Error("parameter author is empty"), http.StatusBadRequest)
// 			return
// 		}
// 		log.Info("query parameter getted", slog.Any("author", author))
// 		quotes, count, err := get.GetWithParam(author)
// 		if err != nil {
// 			log.Error("internal server error", sl.Err(err))
// 			res.Json(w, res.Error("internak server error"), http.StatusInternalServerError)
// 			return
// 		}
// 		log.Info("slice of quotes getted")
// 		respose := GetWithParamResponse{
// 			Quotes: quotes,
// 			Count:  count,
// 		}
// 		log.Info("response id ready", slog.Any("quotes", respose), slog.Int("count", count))
// 		res.Json(w, respose, http.StatusOK)
// 	}
// }
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
