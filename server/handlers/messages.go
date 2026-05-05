package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"relay/models"
)

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")

	h.mu.RLock()
	msgs := make([]*models.Message, 0)
	for _, m := range h.Messages {
		if room == "" || m.Room == room {
			msgs = append(msgs, m)
		}
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

type postMessageRequest struct {
	UserID  string `json:"userId"`
	Content string `json:"content"`
	Room    string `json:"room"`
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	var req postMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}
	if req.Room == "" {
		http.Error(w, "room required", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	user, ok := h.Users[req.UserID]
	h.mu.RUnlock()
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	msg := &models.Message{
		ID:        generateID(),
		UserID:    req.UserID,
		Username:  user.Name,
		Content:   req.Content,
		Room:      req.Room,
		Timestamp: time.Now(),
	}

	h.mu.Lock()
	h.Messages = append(h.Messages, msg)
	h.mu.Unlock()

	// Only clients subscribed to this room receive the broadcast.
	wsMsg := models.WSMessage{Type: "message", Message: msg}
	data, _ := json.Marshal(wsMsg)
	h.Hub.BroadcastToRoom(req.Room, data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
