package server

import (
	"encoding/json"
	"fmt"
	"github.com/Rock-liyi/p2pdb-pubsub/chat/internal/entity"
	"github.com/Rock-liyi/p2pdb-pubsub/chat/protocol"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	bufSize        = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	user   *entity.User
	server *Server
	conn   *websocket.Conn
	send   chan *protocol.Message
}

func NewClient(user *entity.User, server *Server, conn *websocket.Conn) {
	client := &Client{
		user:   user,
		server: server,
		conn:   conn,
		send:   make(chan *protocol.Message, bufSize),
	}
	client.server.register <- client

	go client.readPump()
	go client.writePump()
}

func (c *Client) readMessage() (*protocol.Message, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var message protocol.Message

	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (c *Client) writeMessage(message *protocol.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (c *Client) readPump() {
	defer func() {
		c.server.unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		message, err := c.readMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message.Status = protocol.SucceedStatus
		message.Type = protocol.TextType
		message.FromUser = c.user
		if len(message.Id) == 0 {
			message.Id = uuid.NewV4().String()
		}

		fmt.Printf("client recved message: %v\n", message)

		c.server.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.writeMessage(message)
			if err != nil {
				fmt.Printf("write error: %v\n", err)
				return
			}

			n := len(c.send)
			for i := 0; i < n; i++ {
				_ = c.writeMessage(<-c.send)
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if !r.URL.Query().Has("name") || !r.URL.Query().Has("id") {
		_ = conn.Close()
		return
	}

	user := entity.NewUser(r.URL.Query().Get("id"), r.URL.Query().Get("name"))

	NewClient(user, server, conn)

	fmt.Printf("new client: %s\n", user.DisplayName)
}
