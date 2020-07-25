package nyusocket

// Options module
type Options struct {
	Addr string
}

// Authorization ...
type Authorization struct {
	Authorize chan bool
	Client    *NewClient
}

// ClientMessage ...
type ClientMessage struct {
	Message string
	Client  *Client
}

// Unregister ...
type Unregister struct {
	Store interface{}
	Hub   *Hub
	Hash  string
}
