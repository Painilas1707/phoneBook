package save

import (
	"Study/Demo/internal/StructUser"
	resp "Study/Demo/lib/api/response"
	slogSL "Study/Demo/lib/logger/slog"
	"Study/Demo/storage"
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log"
	"log/slog"
	"net/http"
)

type Request struct {
	ContactFIO  string `json:"contact_fio" validate:"required"`
	BirthDate   string `json:"birth_date" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
}

type Response struct {
	resp.Response
	Error string `json:"error,omitempty"`
	ID    int64  `json:"id,omitempty"`
}

type ContactSaver interface {
	SaveContact(contact StructUser.UserPhoneBook) (int64, error)
}

func New(log *slog.Logger, contactSaver ContactSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.contact.save.New"
		log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", slogSL.Err(err))

			render.JSON(w, r, resp.Error("failed to parse request"))
			return

		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", slogSL.Err(err))

			render.JSON(w, r, resp.Error("failed to validate request"))
			return
		}

		contactFIO := req.ContactFIO
		if contactFIO == "" {
			log.Error("contactFIO field is required", slogSL.Err(err))
		}
		BirthDate := req.BirthDate
		if BirthDate == "" {
			log.Error("BirthDate field is required", slogSL.Err(err))
		}
		phoneNumber := req.PhoneNumber
		if phoneNumber == "" {
			log.Error("phoneNumber field is required", slogSL.Err(err))
		}
		Email := req.Email
		if Email == "" {
			log.Error("Email field is required", slogSL.Err(err))
		}

		id, err := contactSaver.SaveContact(StructUser.UserPhoneBook{})
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Info("contact already exists", slogSL.Err(err))

			render.JSON(w, r, resp.Error("contact already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save contact", slogSL.Err(err))

			render.JSON(w, r, resp.Error("failed to save contact"))
			
			return
		}
		log.Info("contact added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
