package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
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

		if err != nil {
			log.Println("ERROR parsing requiest")
			return
		}

		fmt.Println("POST HANDLER ACTION", action.Action)
		if action.Action == "startRound" {
			go StartRound(gameStore)
		}

		io.WriteString(w, "ok")
	}
}

// StartRound starts a new round
func StartRound(gameStore *ClientsGameCollection) {
	roundState = IS_STARTED
	fmt.Println("Start a round")
	time.Sleep(time.Second * (10 + time.Duration(rand.Intn(5))))
	fmt.Println("Sleep end")

	var ids []string

	for id, item := range *gameStore {
		fmt.Printf("start game for %s %d\n", id, item.Card.Number)
		ids = append(ids, id)
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	var i int8
	for i = 0; i < 3; i++ {
		fmt.Println("_COUNTER:", i+1)
		for id := range ids {
			item := (*gameStore)[ids[id]]
			mc := MessageCounter{
				Type: "counter",
				Payload: MessageCounterPayload{
					Counter: (i + 1),
				},
			}
			SendMessageToGameItemCard(item, Message{
				Type:    mc.Type,
				Payload: mc.Payload,
			})
		}
		TriggerUpdateDashboardState(i+1)
		time.Sleep(2 * time.Second)
	}

	for i := range ids[:(len(ids)/2)] {
		item := (*gameStore)[ids[i]]
		fmt.Println("Close card:", item.Card.ID)

		fmt.Println("Send message card:", item.Card.ID)

		item.Card.IsClosed = true
		item.Card.IsWin = false

		for chanIdx, roundChan := range item.roundChans {
			fmt.Println("Send message card:", item.Card.ID, "for ", chanIdx)
			messageID := uuid.NewV4().String()
			roundMessage := &RoundMessage{
				ID:   messageID,
				Type: "round",
				Payload: RoundMessagePayload{
					IsWin: item.Card.IsWin,
					Card:  *item.Card,
				},
			}
			roundChan <- Message{
				Type:    roundMessage.Type,
				Payload: roundMessage.Payload,
			}
		}
	}
	roundState = IS_FINISHED
	TriggerUpdateDashboardState()
}

// SendMessageToGameItemCard sends messages to chanels
func SendMessageToGameItemCard(gameItem *ClientsGame, message Message) {
	for _, connChan := range gameItem.roundChans {
		fmt.Println("SEND TO GAME ITEM CHAN ", gameItem, " connChan: ", connChan)
		connChan <- message
	}
}
