package nyusocket

// Events ...
type Events struct {
	Authorization chan Authorization
	Register      chan *Client
	Unregister    chan Unregister
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
	Client  *Client
}

// Unregister ...
type Unregister struct {
	Client   *Client
	Continue chan interface{}
}

// NewEvents ...
func NewEvents() *Events {
	return &Events{
		Authorization: make(chan Authorization),
		Register:      make(chan *Client),
		Unregister:    make(chan Unregister),
		ClientMessage: make(chan ClientMessage),
	}
}
