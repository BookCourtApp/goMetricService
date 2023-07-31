package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
	response "github.com/wanna-beat-by-bit/goMetricService/internal/pkg/api"
)

type Response struct {
	response.Response
	Message string `json:"message"`
}

type Saver interface {
	SaveData(metric storage.Metric) error
}

func SaveHandler(saver Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			render.JSON(w, r, Response{
				Message:  "Can't decode a request",
				Response: response.Error("error occured while decoding a request"),
			})
			return
		}
		var bucket storage.Metric
		if err := json.Unmarshal(bytes, &bucket); err != nil {
			render.JSON(w, r, Response{
				Message:  "Can't unmarshal a request",
				Response: response.Error("error occured while unmarshalling a body"),
			})
			return
		}
		defer r.Body.Close()

		log.Println("Got request")
		log.Println(bucket)

		if err := saver.SaveData(bucket); err != nil {
			render.JSON(w, r, Response{
				Response: response.Error("Something went wrong"),
			},
			)
			return
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
			Message:  "You request is good",
		},
		)
	}
}
