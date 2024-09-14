package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	id string

	// connection for the certain client
	socket *websocket.Conn

	// channel to receive messages from other users
	receive chan[]byte

	// room the client is in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		c.room.forward <- message
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
