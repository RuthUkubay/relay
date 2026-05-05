package handlers

import (
	"crypto/rand"
	"fmt"
	"sync"

	"relay/models"
	ws "relay/websocket"
)

// Handler holds all shared server state and is the receiver for every HTTP handler.
type Handler struct {
	mu       sync.RWMutex
	Messages []*models.Message
	Users    map[string]*models.User
	Rooms    map[string]*models.Room // keyed by room name
	Hub      *ws.Hub
}

func NewHandler(hub *ws.Hub) *Handler {
	h := &Handler{
		Messages: make([]*models.Message, 0),
		Users:    make(map[string]*models.User),
		Rooms:    make(map[string]*models.Room),
		Hub:      hub,
	}
	for _, name := range []string{"general", "random", "cs1680"} {
		h.Rooms[name] = &models.Room{ID: name, Name: name}
	}
	return h
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
