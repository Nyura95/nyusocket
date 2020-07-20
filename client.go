package nyusocket

import (
	"log"
	"net/http"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub   *Hub
	Send  chan []byte
	Store interface{}

	conn *websocket.Conn
	hash string
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
		case message, ok := <-c.Send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, newClient *NewClient, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// cookie, err := r.Cookie("X-Auth-Token")
		// if err != nil {
		// 	log.Panicln(err)
		// }
		// log.Println(cookie)
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	send := make(chan []byte, 256)
	client := &Client{hub: hub, conn: conn, Send: send, hash: newClient.getHash(), Store: newClient.Store}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}

func closeServeWs(msg string, key string, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))

	t, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	t.Write(NewMessage("Error", key, msg).Send())
	if err := t.Close(); err != nil {
		return
	}

	log.Printf("Disconnected user during login (%s)", string(msg))
	conn.Close()
}
