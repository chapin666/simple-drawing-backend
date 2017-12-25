package main

import (
	"simple-drawing-backend/utils"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

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

func (client Client) run() {
	go client.read()
	go client.write()
}

func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}
