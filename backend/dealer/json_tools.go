package main

import (
	"encoding/json"
	"net/http"
)

// WriteJSON marshals data and writes it to response 
func WriteJSON(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}