package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	response "github.com/wanna-beat-by-bit/goMetricService/internal/pkg/api"
)

type Response struct {
	response.Response
	Message string `json:"message"`
}

type Saver interface {
	SaveData() error
}

func SaveHandler(saver Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request")
		if err := saver.SaveData(); err != nil {
			render.JSON(w, r, Response{
				Response: response.Error("Something went wrong"),
			},
			)
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
			Message:  "You request is good",
		},
		)
	}
}
