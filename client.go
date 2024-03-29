package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			log.Println(msg)
			c.room.forward <- msg
		} else {
			break
		}
	}
}
func (c *client) write() {
	for msg := range c.send {
		log.Println("send")
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
