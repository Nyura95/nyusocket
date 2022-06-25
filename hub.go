package nyusocket

// Hub of the clients
type Hub struct {
	clients map[*Client]bool
	alive   bool
	info    Info

	register   chan *Client
	unregister chan *Client
	message    chan ClientMessage
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan ClientMessage),
	}
}

func (h Hub) GetInfo() Info {
	return h.info
}

// GetClients return all clients connected
func (h *Hub) GetClients() []*Client {
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// SendToAllClients ...
func (h *Hub) SendToAllClients(message Message) error {
	for client := range h.clients {
		client.send <- message.Send()
	}
	return nil
}

func (h *Hub) getOtherClient(c *Client) []*Client {
	clients := make([]*Client, 0, len(h.clients)-1)
	for client := range h.clients {
		if c.hash != client.hash {
			clients = append(clients, client)
		}
	}
	return clients
}

func (h *Hub) run(events *Events) {
	h.alive = true
	for {
		select {
		case client := <-h.register:
			if events.Register != nil {
				events.Register <- client
			}
			h.clients[client] = true
			h.info.add(client.hash)
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				if events.Unregister != nil {
					events.Unregister <- Unregister{
						Store: c.Store,
						Hub:   h,
						Hash:  c.hash,
					}
				}
				delete(h.clients, c)
				close(c.send)
				h.info.del(c.hash)
			}
		case clientMessage := <-h.message:
			if events.ClientMessage != nil {
				events.ClientMessage <- clientMessage
			}
		}
	}
}
