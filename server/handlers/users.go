package handlers

import (
	"encoding/json"
	"net/http"

	"relay/models"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	users := make([]*models.User, 0, len(h.Users))
	for _, u := range h.Users {
		users = append(users, u)
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

type postUserRequest struct {
	Name string `json:"name"`
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var req postUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	user := &models.User{
		ID:   generateID(),
		Name: req.Name,
	}

	h.mu.Lock()
	h.Users[user.ID] = user
	h.mu.Unlock()

	// user_join goes to all connected clients, not scoped to a room.
	wsMsg := models.WSMessage{Type: "user_join", User: user}
	data, _ := json.Marshal(wsMsg)
	h.Hub.BroadcastToAll(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
