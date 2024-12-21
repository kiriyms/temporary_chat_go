package models

import (
	"bytes"
	"errors"
	"log"
	"sync"
	"text/template"

	"github.com/google/uuid"
)

type Hub struct {
	mu sync.RWMutex

	clients map[*Client]bool
	Id      uuid.UUID

	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client

	messages []*Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Start() {
	for {
		select {
		case msg, ok := <-h.Broadcast:
			if !ok {
				log.Printf("HUB: broadcast channel has been closed. Closing hub %v", h.Id)
				for client := range h.clients {
					log.Printf("HUB: closing client %v", client.Id)
					close(client.Send)
					delete(h.clients, client)
				}
				return
			}
			log.Printf("HUB: hub %v getting broadcast: %v", h.Id, msg)
			h.mu.RLock()

			h.messages = append(h.messages, msg)

			byteTemplate, err := GetTemplateBytes("message-card", msg)
			log.Printf("%v", string(byteTemplate))
			if err != nil {
				log.Printf("HUB: failed to convert broadcasted message to bytes: %v", err)
				byteTemplate = []byte(msg.Content)
			}

			for client := range h.clients {
				select {
				case client.Send <- byteTemplate:
				default:
					log.Printf("HUB: client unresponsive. Closing client %v", client.Id)
					close(client.Send)
					delete(h.clients, client)
				}
			}

			h.mu.RUnlock()

		case client := <-h.Register:
			log.Printf("HUB: hub %v getting register of user %v", h.Id, client.Id)
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			log.Printf("HUB: sending message history from hub %v to client %v", h.Id, client.Id)
			for _, msg := range h.messages {
				byteTemplate, err := GetTemplateBytes("message-card", msg)
				log.Printf("%v", string(byteTemplate))
				if err != nil {
					client.Send <- []byte(msg.Content)
					continue
				}

				client.Send <- byteTemplate
			}

		case client := <-h.Unregister:
			log.Printf("HUB: hub %v getting UNregister of user %v", h.Id, client.Id)
			h.mu.Lock()

			h.clients[client] = false
			if _, ok := h.clients[client]; ok {
				log.Printf("HUB: client unregistered. Closing client %v", client.Id)
				close(client.Send)
				delete(h.clients, client)
			}

			h.mu.Unlock()
		}
	}
}

func GetTemplateBytes(name string, data interface{}) ([]byte, error) {
	t := template.New("")
	_, err := t.ParseGlob("views/*.html")
	if err != nil {
		log.Printf("error parsing blob for template bytes: %v", err)
		return nil, err
	}
	_, err = t.ParseGlob("views/components/*.html")
	if err != nil {
		log.Printf("error parsing blob for template bytes: %v", err)
		return nil, err
	}

	tmpl := t.Lookup(name)
	if tmpl == nil {
		err := errors.New("error parsing template for bytes: template with given name not found")
		log.Printf("%v", err)
		return nil, err
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, data)
	if err != nil {
		log.Printf("error executing template for bytes: %v", err)
		return nil, err
	}

	return renderedMessage.Bytes(), nil
}
