package main

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
}
