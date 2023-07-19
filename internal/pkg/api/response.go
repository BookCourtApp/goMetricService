package response

import "net/http"

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
