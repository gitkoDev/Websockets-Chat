package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {
	sync.RWMutex
	/*
	 'clients' map holds all clients in this room
	  it needs to be a map and not a slice
	  because map has build in delete method
	*/
	clients map[*client]bool

	// channel for new joining clients
	join chan *client

	// channel for leaving clients
	leave chan *client

	forward chan[]byte
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]bool),
		join: make(chan *client),
		leave: make(chan *client),
		forward: make(chan []byte),
	}
}

func(r *room) run() {
	for {
		select {
		case client := <-r.join:
			client.id = uuid.New().String()
			r.clients[client] = true
			fmt.Printf("client %s has joined \n", client.id)
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
			fmt.Printf("client %s has left \n", client.id)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <-msg
			}
		}
	}
}


func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("error upgrading to websocket connection", err)
	}

	client := &client{socket: conn, receive: make(chan []byte, messageBufferSize), room: r,}

	r.join <- client

	defer func(){
		r.leave <- client
	}()

	go client.write()
	client.read()
}
