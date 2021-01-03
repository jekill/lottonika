package main

import (
	"io"
	"net/http"
)

type actionRequestData struct {
	Action string `json:"action"`
}

// CretaeHandleCreateAction creates handler for action
func CretaeHandleCreateAction(gameStore ClientsGameCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "ok")
	}
}
