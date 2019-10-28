package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Card data
type Card struct {
	ID     string `json:"id"`
	Number int16  `json:"number"`
}

// Message object
type Message struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// ClientsGame is information about client's game
type ClientsGame struct {
	Card Card
	ws   *websocket.Conn
}

var upgader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var gameStore = make(map[string]ClientsGame, 200)

var outputBus = make(chan Message)
var inputBus = make(chan Message)
var lastNumberChan = make(chan int16, 1)
var lastNumber int16 = 0
var r = mux.NewRouter()

func main() {

	fmt.Println("Ok let's start")
	r.HandleFunc("/cards", handleCardEndpointPOST).Methods("POST")
	r.HandleFunc("/cards/{id:[^/]+}", handleCardEndpointGET)
	r.HandleFunc("/cards/{id}/ws", handleWSConnections)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>OK</body></html>")
	})

	go lastNumberGenerator()
	go watchBus()

	srv := &http.Server{
		Addr: "localhost:5555",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	err := srv.ListenAndServe()
	if err != nil {
		println("We have a problem")
		log.Fatal(err)
	}

}

func lastNumberGenerator() {
	for {
		lastNumber++
		lastNumberChan <- lastNumber
	}
}

func handleCardEndpointGET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*Card create")
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	item, ok := gameStore[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, item.Card)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleCardEndpointPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	uuid, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	card := Card{
		ID:     uuid.String(),
		Number: <-lastNumberChan,
	}

	gameStore[card.ID] = ClientsGame{
		Card: card,
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, card)
}

func handleWSConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start WS")
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println(id)
	item, ok := gameStore[id]
	fmt.Println(item.Card.ID)

	if !ok {
		fmt.Println("Id not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ws, err := upgader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("Error of message")
			continue
		}

		outMsg := <-outputBus
		ws.WriteJSON(outMsg)
	}
}

func watchBus() {
	for {
		time.Sleep(time.Second * 3)

		outMsg := Message{
			ID:      "1",
			Type:    "hello out",
			Payload: "none",
		}

		outputBus <- outMsg

	}
}
