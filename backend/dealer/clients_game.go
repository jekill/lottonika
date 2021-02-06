package main

import (
	"fmt"
	"sync"
)

// ClientsGame is information about client's game
type ClientsGame struct {
	Card     *Card
	mtx      sync.Mutex
	conChans []chan Message
}

func (c *ClientsGame) addConChan() (index int, conChan chan Message) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	conChan = make(chan Message, 100)
	c.conChans = append(c.conChans, conChan)
	fmt.Println("Added a new card con chan", c.conChans)

	return len(c.conChans) - 1, conChan
}

func (c *ClientsGame) removeConChan(index int) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	println("REmove a chan withindex ", index, "len", len(c.conChans), "channs:", c.conChans)
	conChan := c.conChans[index]
	close(conChan)
	c.conChans = append(c.conChans[:index], c.conChans[(index+1):]...)
}

func (c *ClientsGame) sendMessage(message Message) {
	for _, connChan := range c.conChans {
		fmt.Println("SEND TO GAME ITEM CHAN ", c, " connChan: ", connChan)
		connChan <- message
	}
}
