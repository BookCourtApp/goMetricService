package saveHandler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
	response "github.com/wanna-beat-by-bit/goMetricService/internal/pkg/response"
)

type Response struct {
	response.Response
	Message string `json:"message"`
}

type LogSaver interface {
	SaveData(metric storage.Metric) error
}

func New(saver LogSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var bucket storage.Metric
		if err := render.DecodeJSON(r.Body, &bucket); err != nil {
			render.Render(w, r, response.ErrBadRequest(err))
		}
		defer r.Body.Close()

		if err := saver.SaveData(bucket); err != nil {
			render.Render(w, r, response.ErrInternalServerError(err))
			return
		}

		render.Render(w, r, response.Created())
	}
}
