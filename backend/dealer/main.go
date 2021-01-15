package main

import (
	// "encoding/json"
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
	ID       string `json:"id"`
	Number   int16  `json:"number"`
	IsClosed bool   `json:"is_closed"`
	IsWin    bool   `json:"is_win"`
}

// Message object
type Message struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// RoundMessage a message for UI
type RoundMessage struct {
	ID      string              `json:"id"`
	Type    string              `json:"type"`
	Payload RoundMessagePayload `json:"payload"`
}

// RoundMessagePayload is the payload
type RoundMessagePayload struct {
	IsWin bool `json:"is_win"`
	Card  Card `json:"card"`
}

// IsWinMessage is the message about resulf of user's round
type IsWinMessage struct {
	IsWin bool `json:"isWin"`
}

// ClientsGame is information about client's game
type ClientsGame struct {
	Card *Card
	ws   *websocket.Conn
	// roundChan chan RoundMessage
	roundChans []chan IsWinMessage
}

// ClientsGameCollection store or collection map "card id" -> ClientGame{}
type ClientsGameCollection map[string]*ClientsGame

var upgader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var gameStore = make(ClientsGameCollection, 200)

var outputBus = make(chan Message)
var inputBus = make(chan Message)

// number on the card autoincremented
var lastNumberChan = make(chan int16, 1)
var lastNumber int16 = 0
var r = mux.NewRouter()

func main() {

	fmt.Println("Ok let's start")
	r.HandleFunc("/cards", handleCardEndpointPOST).Methods("POST")
	r.HandleFunc("/cards/{id:[^/]+}", handleCardEndpointGET)
	r.HandleFunc("/cards/{id}/ws", handleWSConnections)

	r.HandleFunc("/game-manager/actions", CretaeHandleCreateAction(&gameStore)).Methods(http.MethodPost)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>OK</body></html>")
	})

	go lastNumberGenerator()

	// go startRound()

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
	WriteJSON(w, item.Card)
}

func handleCardEndpointPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	uuid := uuid.NewV4()

	card := Card{
		ID:     uuid.String(),
		Number: <-lastNumberChan,
	}

	gameStore[card.ID] = &ClientsGame{
		Card:      &card,
		roundChans: make([]chan IsWinMessage, 1),
	}

	w.WriteHeader(http.StatusCreated)
	WriteJSON(w, card)
}

func handleWSConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start WS")
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Handle ws connect", id)
	item, ok := gameStore[id]
	fmt.Println("Item card is found in the storage", item.Card.ID)

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
	roundChan := make(chan IsWinMessage, 1)

	item.roundChans = append(item.roundChans, roundChan)
	roundChanIndex := len(item.roundChans) - 1

	defer func() {
		ws.Close()
		fmt.Println("Websocket closed for card:", item.Card.ID)
		item.roundChans = append(item.roundChans[:roundChanIndex], item.roundChans[(roundChanIndex+1):]...)
	}()

	// item.ws = ws

	ws.SetCloseHandler(func(code int, text string) error {
		fmt.Println("WSocket closeHandler:", item.Card.ID)
		item.ws = nil
		return nil
	})

	for {
		<-roundChan

		messageID := uuid.NewV4().String()

		ws.WriteJSON(&RoundMessage{
			ID:   messageID,
			Type: "round",
			Payload: RoundMessagePayload{
				IsWin: item.Card.IsWin,
				Card:  *item.Card,
			},
		})
	}
}

