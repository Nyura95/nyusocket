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
	Clients map[*Client]bool

	Broadcast               chan []byte
	disseminateToTheOthers  chan *DisseminateToTheOthers
	DisseminateToTheTargets chan *DisseminateToTheTargets

	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:               make(chan []byte),
		DisseminateToTheTargets: make(chan *DisseminateToTheTargets),

		disseminateToTheOthers: make(chan *DisseminateToTheOthers),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run(events *Events) {
	for {
		select {
		case client := <-h.register:
			if events.Register != nil {
				events.Register <- client
			}
			h.Clients[client] = true
			Infos.Add(client.Hash)

		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				Infos.Del(client.Hash)
				if events.Unregister != nil {
					events.Unregister <- h.Clients
				}
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					h.unregister <- client
				}
			}
		}
	}
}
