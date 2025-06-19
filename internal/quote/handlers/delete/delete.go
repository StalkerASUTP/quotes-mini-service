package delete

import (
	"log/slog"
	"net/http"
	"quotes-mini-service/pkg/res"
	"quotes-mini-service/pkg/sl"
	"strconv"
	"strings"
)

type QuoteDeleter interface {
	Delete(id int) error
}

func New(log *slog.Logger, delete QuoteDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"
		log = log.With(
			slog.String("op", op),
		)
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid argument", sl.Err(err))
			res.Json(w, res.Error("invalid argument"), http.StatusBadRequest)
			return
		}
		log.Info("id converted", slog.Any("id", id))
		err = delete.Delete(id)
		if err != nil {
			log.Error(err.Error(), sl.Err(err))
			if strings.Contains(err.Error(), "not found") {
				res.Json(w, res.Error("entry with this id not found"), http.StatusNotFound)
			} else {
				res.Json(w, res.Error("internal server error"), http.StatusInternalServerError)
			}
			return
		}
		log.Info("quote deleted", slog.Int64("id", int64(id)))
		res.Json(w, nil, http.StatusNoContent)
	}
}
