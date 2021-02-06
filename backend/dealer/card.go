package main

// Card data
type Card struct {
	ID       string `json:"id"`
	Number   uint16 `json:"number"`
	IsClosed bool   `json:"is_closed"`
	IsWin    bool   `json:"is_win"`
}


