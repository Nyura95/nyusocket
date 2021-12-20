package nyusocket

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub   *Hub
	send  chan []byte
	Store *Store

	query map[string][]string
	path  string

	conn *websocket.Conn
	hash string
}

// GetPath return the path client
func (c *Client) GetPath() string {
	return c.path
}

// GetQuery return query client
func (c *Client) GetQuery(key string) ([]string, error) {
	if val, exist := c.query[key]; exist {
		return val, nil
	}
	return nil, errors.New("key not exist")
}

// Close disconnect the client
func (c *Client) Close() {
	c.hub.unregister <- c
}

// Send a message
func (c *Client) Send(message []byte) error {
	if !Infos.Alive(c) {
		return errors.New("client unregisted")
	}
	c.send <- message
	return nil
}

// GetOthersClients ...
func (c *Client) GetOthersClients() []*Client {
	return c.hub.getOtherClient(c)
}

// GetHash ...
func (c *Client) GetHash() string {
	return c.hash
}

// GetAllClients ...
func (c *Client) GetAllClients() []*Client {
	return c.hub.GetClients()
}

func (c *Client) getHash() string {
	return c.hash
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.hub.message <- ClientMessage{
			Message: string(message),
			Client:  c,
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}

			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					return
				}
				_, err = w.Write(<-c.send)
				if err != nil {
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
