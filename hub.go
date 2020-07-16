package nyusocket

// Hub of the clients
type Hub struct {
	clients map[*Client]bool

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

// GetClient return all clients connected
func (h *Hub) GetClient() []*Client {
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// GetOtherClient return all clients connected without hash
func (h *Hub) GetOtherClient(c *Client) []*Client {
	clients := make([]*Client, 0, len(h.clients)-1)
	for client := range h.clients {
		if c.hash != client.hash {
			clients = append(clients, client)
		}
	}
	return clients
}

func (h *Hub) run(events *Events) {
	for {
		select {
		case client := <-h.register:
			if events.Register != nil {
				events.Register <- client
			}
			h.clients[client] = true
			Infos.add(client.hash)
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				if events.Unregister != nil {
					trigger := make(chan interface{})
					events.Unregister <- Unregister{
						Client:   c,
						Continue: trigger,
					}
					<-trigger
					close(trigger)
				}
				delete(h.clients, c)
				close(c.Send)
				Infos.del(c.hash)
			}
		case clientMessage := <-h.message:
			events.ClientMessage <- clientMessage
		}
	}
}
