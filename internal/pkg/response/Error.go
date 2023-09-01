package response

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	StatusBadRequest          = "Invalid request"
	StatusInternalServerError = "Internal sever error"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: http.StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: http.StatusInternalServerError,
		Error:  msg,
	}
}

type ErrResponse struct {
	Error          error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	ErrorText  string `json:"error"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Error:          err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     StatusBadRequest,
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Error:          err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     StatusInternalServerError,
		ErrorText:      err.Error(),
	}
}
