package main

import (
	"encoding/json"
	"log"
	"net/http"
	"simple-drawing-backend/message"

	"github.com/tidwall/gjson"

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
	log.Println("client connected: ", client.socket.RemoteAddr())

	// Make list of all users
	users := []message.User{}
	for _, c := range hub.clients {
		users = append(users, message.User{ID: c.id, Color: c.color})
	}

	// Send exists clients to client
	hub.send(message.NewConnected(client.color, users), client)
	// Notify user joined
	hub.broadcast(message.NewUserJoined(client.id, client.color), client)
}

// onDisconnect method of Hub.
func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())
	client.close()

	// Find index of client
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}

	// Delete client from list
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

	// Notify user left
	hub.broadcast(message.NewUserLeft(client.id), nil)
}

// onMessage method of Hub.
func (hub *Hub) onMessage(data []byte, client *Client) {
	kind := gjson.GetBytes(data, "kind").Int()
	if kind == message.KindStroke {
		var msg message.Stroke
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = client.id
		hub.broadcast(msg, client)
	} else if kind == message.KindClear {
		var msg message.Clear
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = client.id
		hub.broadcast(msg, client)
	}
}
