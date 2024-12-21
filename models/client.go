package models

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   uuid.UUID
	Hub  *Hub
	conn *websocket.Conn
	Send chan []byte
}

const writeWait = 5 * time.Second
const pongWait = 10 * time.Second
const pingPeriod = (pongWait * 9) / 10
const maxMsgSize = 512

func NewClient(id uuid.UUID, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		Id:   id,
		Hub:  hub,
		conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (c *Client) WriteToWebSocket() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Printf("CLIENT: closing WRITE goroutine, stopping ticker, closing websocket of client %v", c.Id)
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			log.Printf("CLIENT: client %v received msg: %v", c.Id, string(msg))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Printf("CLIENT: client %v SEND channel closed. Writing CloseMessage", c.Id)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			err = w.Close()
			if err != nil {
				return
			}

		case <-ticker.C:
			log.Printf("CLIENT: pinging to client %v", c.Id)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadFromWebSocket() {
	defer func() {
		log.Printf("CLIENT: closing READ goroutine, unregistering client %v from hub %v", c.Id, c.Hub.Id)
		c.Hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.conn.ReadMessage()
		log.Printf("CLIENT: received from ws: %v", string(text))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("CLIENT: unexpected close error: %v", err)
			}
			break
		}

		msg := &Message{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			log.Println("CLIENT: failed to decode: ", err)
		}

		msg.Id = c.Id
		log.Printf("CLIENT: decoded message info: id: %v, text: %v", msg.Id, msg.Content)
		c.Hub.Broadcast <- msg
	}
}
