package models

import "net/http"

type HttpResponse struct {
	Message    string         `json:"message,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

func EnvelopIntoHttpResponse(obj any, objName string, code int) HttpResponse {
	message := http.StatusText(code)
	return HttpResponse{
		Message: message,
		Properties: map[string]any{
			objName: obj,
		},
	}
}
