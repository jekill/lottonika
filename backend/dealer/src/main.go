package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Message object
type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

var upgader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool, 200)

func main() {

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		println("We have a problem")
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

}
