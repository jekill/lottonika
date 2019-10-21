package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message object
type Message struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

var upgader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool, 200)
var outputBus = make(chan Message)
var inputBus = make(chan Message)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>OK</body></html>")
	})

	http.HandleFunc("/lott", handleWSConnections)

	err := http.ListenAndServe("localhost:5555", nil)
	if err != nil {
		println("We have a problem")
		log.Fatal(err)
	}

	go watchBus()
}

func handleWSConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		// fmt.Fprintf(w, "Error")
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("Error of message")
			delete(clients, ws)
			break
		}

		inputBus <- msg

		outMsg := <-outputBus
		ws.WriteJSON(outMsg)
	}
}

func watchBus() {
	for {
		<-inputBus

		outMsg := Message{
			ID:      "1",
			Type:    "hello out",
			Payload: "none",
		}

		outputBus <- outMsg

	}
}
