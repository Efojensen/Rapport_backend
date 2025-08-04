package models

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type WSMessage struct {
	Type   string      `json:"type"`
	RoomId string      `json:"roomId"`
	UserId string      `json:"userId"`
	Data   interface{} `json:"data"`
	Chat   *Chat       `json:"chat,omitempty"`
}

type Client struct {
	Conn   *websocket.Conn
	UserId string
	RoomId string
	Send   chan WSMessage
}

type Hub struct {
	Clients    map[string]*Client
	Rooms      map[string]map[string]*Client // roomId -> userId -> Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan WSMessage
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan WSMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastToRoom(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.Clients[client.UserId] = client

	if h.Rooms[client.RoomId] == nil {
		h.Rooms[client.RoomId] = make(map[string]*Client)
	}
	h.Rooms[client.RoomId][client.UserId] = client

	log.Printf("Client %s joined room %s", client.UserId, client.RoomId)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.Clients[client.UserId]; ok {
		delete(h.Clients, client.UserId)
		close(client.Send)
	}

	if room, ok := h.Rooms[client.RoomId]; ok {
		if _, ok := room[client.UserId]; ok {
			delete(room, client.UserId)
			if len(room) == 0 {
				delete(h.Rooms, client.RoomId)
			}
		}
	}

	log.Printf("Client %s left room %s", client.UserId, client.RoomId)
}

func (h *Hub) broadcastToRoom(message WSMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if room, ok := h.Rooms[message.RoomId]; ok {
		for _, client := range room {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client.UserId)
				delete(room, client.UserId)
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error writing message: %v", err)
			return
		}
	}

	// Channel closed, send close message
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		message.UserId = c.UserId
		message.RoomId = c.RoomId
		hub.Broadcast <- message
	}
}
