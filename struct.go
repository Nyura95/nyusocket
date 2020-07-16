package nyusocket

// Options module
type Options struct {
	Port int
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
	Client   *Client
	Continue chan interface{}
}
