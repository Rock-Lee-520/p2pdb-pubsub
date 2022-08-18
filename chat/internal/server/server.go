package server

import (
	"fmt"
	"github.com/Rock-liyi/p2pdb-pubsub/chat/protocol"
)

type Server struct {
	clients    map[*Client]bool
	broadcast  chan *protocol.Message
	register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *protocol.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Server) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
		case message := <-h.broadcast:
			for client := range h.clients {
				if client.user.Id == message.FromUser.Id {
					continue
				}
				fmt.Printf("send message [%s] to user: [%s]\n", message.Content, client.user.DisplayName)
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}
