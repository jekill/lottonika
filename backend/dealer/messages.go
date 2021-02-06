package main

const (
	MessageTypeFullState = "state"
	MessageTypeCounter   = "counter"
	MessageTypeRound     = "round"
	MessageTypeCardState = "card_state"
)

// Message object
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// StateMessagePayload a struct for payload
type StateMessagePayload struct {
	Cards        []Card `json:"cards,omitempty"`
	RoundState   string `json:"round_state,omitempty"`
	Counter      string `json:"counter,omitempty"`
	CurrentRound uint8  `json:"current_round,omitempty"`
}

// MessageCounterPayload struct
type MessageCounterPayload struct {
	Counter string `json:"counter,omitempty"`
}

type MessageCardStatePayload struct {
	Card  Card `json:"card,omitempty"`
	RoundState   string `json:"round_state,omitempty"`
	Counter      string `json:"counter,omitempty"`
	CurrentRound uint8  `json:"current_round,omitempty"`
}

// RoundMessagePayload is the payload
type RoundMessagePayload struct {
	IsWin bool `json:"is_win"`
	Card  Card `json:"card"`
}
