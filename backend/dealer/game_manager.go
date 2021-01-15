package main

import (
	"math/rand"
	"time"
	"fmt"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type actionRequestData struct {
	Action string `json:"action"`
}

// CretaeHandleCreateAction creates handler for action
func CretaeHandleCreateAction(gameStore *ClientsGameCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("POST HANDLER")
		decoder := json.NewDecoder(r.Body)
		var action actionRequestData
		err := decoder.Decode(&action)

		if err!=nil{
			log.Println("ERROR parsing requiest")
			return;
		}

		fmt.Println("POST HANDLER ACTION", action.Action)
		if action.Action=="startRound"{
			go StartRound(gameStore)
		}

		io.WriteString(w, "ok")
	}
}


// StartRound starts a new round
func StartRound(gameStore *ClientsGameCollection) {
	fmt.Println("Start a round")
	time.Sleep(time.Second * (10 + time.Duration(rand.Intn(20))))
	fmt.Println("Sleep end")

	var ids []string

	for id, card := range *gameStore {
		fmt.Printf("start game for %s %d\n", id, card.Card.Number)
		ids = append(ids, id)
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	for i := range ids {
		item := (*gameStore)[ids[i]]
		fmt.Println("Close card:", item.Card.ID)

		fmt.Println("Send message card:", item.Card.ID)

		item.Card.IsClosed = true
		item.Card.IsWin = false

		for chanIdx, roundChan := range item.roundChans {
			fmt.Println("Send message card:", item.Card.ID, "for ", chanIdx)
			roundChan <- IsWinMessage{
				IsWin: item.Card.IsWin,
			}
		}
	}
}
