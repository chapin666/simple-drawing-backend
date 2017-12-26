package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub struct.
type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
}

// newHub function create new Hub struct.
func newHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// run method of Hub
func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

// handleWebSocket method process socket message.
func (hub *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// check client is supported for websocket.
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.clients = append(hub.clients, client)
	hub.register <- client
	client.run()
}

// send method.
func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

// broadcast method broadcasts a message to all clients, except one(sender).
func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

// onConnect method of Hub.
func (hub *Hub) onConnect(client *Client) {

}

// onDisconnect method of Hub.
func (hub *Hub) onDisconnect(client *Client) {

}

// onMessage method of Hub.
func (hub *Hub) onMessage(data []byte, client *Client) {

}
