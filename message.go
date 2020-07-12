package nyusocket

import (
	"encoding/json"
	"time"
)

// Message action front
type Message struct {
	Action  string
	Message string
	Key     string
	Created time.Time
}

// NewMessage create a new instance of message
func NewMessage(action, message, key string) *Message {
	return &Message{Action: action, Message: message, Key: key, Created: time.Time.UTC(time.Now())}
}

// Send to the front
func (n *Message) Send() []byte {
	json, _ := json.Marshal(n)
	return append(json)
}
