package save

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

const aliasLength = 6 //TODO: move to config

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSavergo:generate
type URLSaver interface {
	SaveUrl(alias string, urlToStore string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "emy request",
			})
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			// convert the error to a validation error
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			//render.JSON(w, r, "internal server error")
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err = urlSaver.SaveUrl(alias, req.URL)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exist", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exist"))
			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added")

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
