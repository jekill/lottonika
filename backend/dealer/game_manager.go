package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
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
	currentRound++
	updateRoundState(RoundStateStarted)
	fmt.Println("Start a round")
	time.Sleep(time.Second * 5)
	//time.Sleep(time.Second * (10 + time.Duration(rand.Intn(5))))
	fmt.Println("Sleep end")

	var ids []string

	for id, item := range *gameStore {
		fmt.Printf("start game for %s %d\n", id, item.Card.Number)
		ids = append(ids, id)
	}

	cardsLen := len(ids)
	if cardsLen < 2 {
		updateRoundState(RoundStateFinished)
		TriggerUpdateDashboardState()
		return
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	for _, i := range []string{"1", "2", "2.5", "2.9", "🤪", "2.999", "3"} {
		fmt.Println("_COUNTER:", i)
		for id := range ids {
			item := (*gameStore)[ids[id]]
			item.sendMessage(Message{
				Type: MessageTypeCounter,
				Payload: MessageCounterPayload{
					Counter: i,
				},
			})
		}
		TriggerUpdateDashboardState(i)
		time.Sleep(2 * time.Second)
	}

	maxFailedIdx := cardsLen / 2
	for i := range ids {
		item := (*gameStore)[ids[i]]
		fmt.Println("Close card:", item.Card.ID)

		fmt.Println("Send message card:", item.Card.ID)

		if i <= maxFailedIdx {
			item.Card.IsClosed = true
			item.Card.IsWin = false
		} else {
			item.Card.IsClosed = false
			item.Card.IsWin = true
		}

		item.sendMessage(Message{
			Type: MessageTypeCardState,
			Payload: MessageCardStatePayload{
				Card:         *item.Card,
				RoundState:   roundState,
				CurrentRound: currentRound,
			},
		})
	}
	updateRoundState(RoundStateFinished)

	TriggerUpdateDashboardState()
	time.Sleep(5 * time.Second)

	RemoveFailedCardsAndEmptyWins()

	updateRoundState(RoundStateStarted)
}

func updateRoundState(state string) {
	roundState = state
	TriggerUpdateDashboardState()
	for _, item := range gameStore {
		item.sendMessage(Message{
			Type: MessageTypeCardState,
			Payload: MessageCardStatePayload{
				Card: *item.Card,
				CurrentRound: currentRound,
				RoundState:   roundState,
			},
		})
	}
}

func RemoveFailedCardsAndEmptyWins() {
	for id, item := range gameStore {
		if item.Card.IsClosed && !item.Card.IsWin {
			delete(gameStore, id)
		} else {
			item.Card.IsWin = false
		}
	}
}
