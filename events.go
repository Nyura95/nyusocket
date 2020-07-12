package nyusocket

// Events ...
type Events struct {
	Authorization chan Authorization
	Register      chan *Client
	Unregister    chan map[*Client]bool
}

// Authorization ...
type Authorization struct {
	Hash      string
	Authorize chan bool
}

// NewEvents ...
func NewEvents() *Events {
	return &Events{
		Authorization: make(chan Authorization),
		Register:      make(chan *Client),
		Unregister:    make(chan map[*Client]bool),
	}
}
