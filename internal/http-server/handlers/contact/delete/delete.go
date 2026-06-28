package delete

import (
	resp "Study/Demo/lib/api/response"
	slogSL "Study/Demo/lib/logger/slog"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type ContactDeletter interface {
	DeleteContact(id string) error
}

func New(log *slog.Logger, contactDelete ContactDeletter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.contact.delete.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("Контакт с таким ID не существует")

			render.JSON(w, r, resp.Error("Невозможно найти строку, как как строка пустая"))

			return
		}

		err := contactDelete.DeleteContact(id)
		if err != nil {
			log.Error("Не удалось удалить контакт", slogSL.Err(err))
			render.JSON(w, r, resp.Error("Не найдено"))

			return
		}
		log.Info("Контакт удален", slog.String("id", id))

		render.JSON(w, r, resp.OK())
	}

}
