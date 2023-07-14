package response

import "net/http"

type Response struct {
	Status int   `json:"status"`
	Error  error `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: http.StatusOK,
	}
}
