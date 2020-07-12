package nyusocket

// Events ...
type Events struct {
	Authorization chan Authorization
	Register      chan *Client
	Unregister    chan map[*Client]bool
	ClientMessage chan ClientMessage
}

// Authorization ...
type Authorization struct {
	Hash      string
	Authorize chan bool
}

// ClientMessage ...
type ClientMessage struct {
	Message string
	client  *Client
}

// NewEvents ...
func NewEvents() *Events {
	return &Events{
		Authorization: make(chan Authorization),
		Register:      make(chan *Client),
		Unregister:    make(chan map[*Client]bool),
		ClientMessage: make(chan ClientMessage),
	}
}
