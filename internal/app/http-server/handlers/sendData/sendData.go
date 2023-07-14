package sendData

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

func New() http.HandlerFunc {
	log.Println("Endpoint 'sendData' is up")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request")

		render.JSON(w, r, Response{
			Response: response.OK(),
			Message:  "You request is good",
		},
		)
	}
}
