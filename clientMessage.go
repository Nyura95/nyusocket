package nyusocket

import "encoding/json"

type ClientMessage struct {
	Message string
	Client  *Client
}

func (c ClientMessage) ParseMessage(data interface{}) error {
	return json.Unmarshal([]byte(c.Message), &data)
}
