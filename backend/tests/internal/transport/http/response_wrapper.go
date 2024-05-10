package http

import (
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api"
	"net/http"
)

func WriteError(w http.ResponseWriter, e error) {
	api.WriteError(w, CodeFromError(e), ErrorCode(e), e.Error())
}

func WriteErrorMessage(w http.ResponseWriter, e error, msg string) {
	api.WriteErrorMessage(w, CodeFromError(e), ErrorCode(e), msg)
}

func WriteResponse(w http.ResponseWriter, code int, payload interface{}) {
	api.WriteResponse(w, code, payload)
}

func WriteErrorManual(w http.ResponseWriter, code int, e api.Error) {
	api.WriteErrorManual(w, code, e)
}
