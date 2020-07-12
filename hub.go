package nyusocket

// DisseminateToTheOthers for send a message to all account except the broadcaster
type DisseminateToTheOthers struct {
	broadcaster int
	message     []byte
}

// DisseminateToTheTargets for send message to the target
type DisseminateToTheTargets struct {
	Targets []int
	Message []byte
}

// Hub of the clients
type Hub struct {
	clients map[*Client]bool

	Broadcast               chan []byte
	disseminateToTheOthers  chan *DisseminateToTheOthers
	DisseminateToTheTargets chan *DisseminateToTheTargets

	register   chan *Client
	unregister chan *Client
	message    chan ClientMessage
}

func newHub() *Hub {
	return &Hub{
		Broadcast:               make(chan []byte),
		DisseminateToTheTargets: make(chan *DisseminateToTheTargets),

		disseminateToTheOthers: make(chan *DisseminateToTheOthers),
		clients:                make(map[*Client]bool),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan ClientMessage),
	}
}

// GetClient return all clients connected
func (h *Hub) GetClient() []*Client {
	clients := make([]*Client, len(h.clients), 0)
	for client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// GetOtherClient return all clients connected without hash
func (h *Hub) GetOtherClient(hash string) []*Client {
	clients := make([]*Client, len(h.clients)-1, 0)
	for client := range h.clients {
		if hash != client.Hash {
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
			Infos.Add(client.Hash)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				Infos.Del(client.Hash)
				if events.Unregister != nil {
					events.Unregister <- h.clients
				}
			}
		case clientMessage := <-h.message:
			events.ClientMessage <- clientMessage
		}
	}
}
