package models

import "time"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Room      string    `json:"room"`
	Timestamp time.Time `json:"timestamp"`
}

// WSMessage is what the hub delivers to every relevant WebSocket client.
type WSMessage struct {
	Type    string   `json:"type"` // "message" | "user_join"
	Message *Message `json:"message,omitempty"`
	User    *User    `json:"user,omitempty"`
}
