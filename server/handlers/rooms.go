package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"relay/models"
)

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	rooms := make([]*models.Room, 0, len(h.Rooms))
	for _, room := range h.Rooms {
		rooms = append(rooms, room)
	}
	h.mu.RUnlock()

	sort.Slice(rooms, func(i, j int) bool { return rooms[i].Name < rooms[j].Name })

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}
