package main

import (
	"net/http"
)

func DisplayError(message string, errorCode int, w http.ResponseWriter) {
	w.WriteHeader(errorCode)
	w.Write([]byte(`{ "message": "` + message + `"}`))
}

