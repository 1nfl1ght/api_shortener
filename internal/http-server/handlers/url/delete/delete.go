package delete

import (
	resp "api-shorter/internal/lib/api/response"
	"api-shorter/internal/lib/logger/sl"
	"api-shorter/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log.Info("ALO BLYAT")

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		err := urlDeleter.DeleteURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted url with alias")
	}
}
