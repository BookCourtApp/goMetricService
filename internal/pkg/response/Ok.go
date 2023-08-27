package response

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	StatusOK      = "OK"
	StatusCreated = "Created"
)

type OkResponse struct {
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
}

func (ok *OkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, ok.HTTPStatusCode)
	return nil
}

func Ok() render.Renderer {
	return &OkResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText:     StatusOK,
	}
}

func Created() render.Renderer {
	return &OkResponse{
		HTTPStatusCode: http.StatusCreated,
		StatusText:     StatusCreated,
	}
}
