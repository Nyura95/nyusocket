package nyusocket

import "github.com/segmentio/ksuid"

// Options module
type Options struct {
	Addr string
}

// BeforeUpgrade ...
type BeforeUpgrade struct {
	Authorize chan bool
	Client    *Client
}

// ClientMessage ...
type ClientMessage struct {
	Message string
	Client  *Client
}

// Unregister ...
type Unregister struct {
	Store *Store
	Hub   *Hub
	Hash  string
}

// NewClient ...
type NewClient struct {
	Query map[string][]string
	Path  string
	Store *Store
	Hash  string
}

func (n *NewClient) getHash() string {
	if n.Hash == "" {
		n.Hash = ksuid.New().String()
	}
	return n.Hash
}
