package socket

import "github.com/segmentio/ksuid"

// NewClient ...
type NewClient struct {
	Query map[string][]string
	Store interface{}
	Hash  string
}

func (n *NewClient) getHash() string {
	if n.Hash == "" {
		n.Hash = ksuid.New().String()
	}
	return n.Hash
}
