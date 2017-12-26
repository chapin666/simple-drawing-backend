package main

import (
	"simple-drawing-backend/utils"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Client struct.
type Client struct {
	id       string
	hub      *Hub
	color    string
	socket   *websocket.Conn
	outbound chan []byte
}

// newClient return a new client
func newClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		id:       uuid.NewV4().String(),
		color:    utils.GenerateColor(),
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

// read method.
func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
	}()

	for {
		_, data, err := client.socket.ReadMessage()
		if err != nil {
			break
		}
		client.hub.onMessage(data, client)
	}
}

// write method.
func (client *Client) write() {
	for {
		select {
		case data, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// run method
func (client Client) run() {
	go client.read()
	go client.write()
}

// close method
func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}
