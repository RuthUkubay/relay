package websocket

import (
	"log"

	gws "github.com/gorilla/websocket"
)

// BroadcastMsg pairs a payload with an optional room scope.
// An empty Room sends to every connected client.
type BroadcastMsg struct {
	Room string
	Data []byte
}

// Client is a single WebSocket connection paired with its outbound send buffer.
type Client struct {
	Hub    *Hub
	Conn   *gws.Conn
	Send   chan []byte
	UserID string
	Room   string
}

// Hub owns the set of active clients and fans out broadcasts.
// All mutations to clients happen inside Run(), so no mutex is needed on the map.
type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan BroadcastMsg
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan BroadcastMsg, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// BroadcastToRoom sends data only to clients currently in the given room.
func (h *Hub) BroadcastToRoom(room string, data []byte) {
	h.Broadcast <- BroadcastMsg{Room: room, Data: data}
}

// BroadcastToAll sends data to every connected client regardless of room.
func (h *Hub) BroadcastToAll(data []byte) {
	h.Broadcast <- BroadcastMsg{Data: data}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}

		case msg := <-h.Broadcast:
			for client := range h.clients {
				// Skip clients in other rooms when the broadcast is room-scoped.
				if msg.Room != "" && client.Room != msg.Room {
					continue
				}
				select {
				case client.Send <- msg.Data:
				default:
					delete(h.clients, client)
					close(client.Send)
				}
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if gws.IsUnexpectedCloseError(err, gws.CloseGoingAway, gws.CloseAbnormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(gws.TextMessage, msg); err != nil {
			return
		}
	}
}
