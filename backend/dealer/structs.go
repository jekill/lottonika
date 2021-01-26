package main

var MESSAGE_TYPE_FULL_STATE = "state"

// StateMessage a struct for a message
type StateMessage struct {
	Type    string              `json:"type"`
	Payload StateMessagePayload `json:"payload"`
}

// StateMessagePayload a struct for payload
type StateMessagePayload struct {
	Cards      []Card `json:"cards"`
	RoundState int    `json:"round_state"`
	Counter    int8   `json:"counter,omitempty"`
}
