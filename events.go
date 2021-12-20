package nyusocket

import "log"

// Events ...
type Events struct {
	BeforeUpgrade chan BeforeUpgrade
	Register      chan *Client
	Unregister    chan Unregister
	ClientMessage chan ClientMessage
}

// NewEvents ...
func NewEvents() *Events {
	return &Events{}
}

// CreateAllEvents ...
func (e *Events) CreateAllEvents() {
	e.CreateBeforeUpgradeEvent()
	e.CreateClientMessageEvent()
	e.CreateRegisterEvent()
	e.CreateUnregisterEvent()
}

// CreateBeforeUpgradeEvent ...
func (e *Events) CreateBeforeUpgradeEvent() {
	e.BeforeUpgrade = make(chan BeforeUpgrade)
}

// CreateRegisterEvent ...
func (e *Events) CreateRegisterEvent() {
	e.Register = make(chan *Client)
}

// CreateUnregisterEvent ...
func (e *Events) CreateUnregisterEvent() {
	e.Unregister = make(chan Unregister)
}

// CreateClientMessageEvent ...
func (e *Events) CreateClientMessageEvent() {
	e.ClientMessage = make(chan ClientMessage)
}

// Close all chan event
func (e *Events) Close() {
	log.Println(e.BeforeUpgrade)
	if e.BeforeUpgrade != nil {
		close(e.BeforeUpgrade)
	}
	if e.Register != nil {
		close(e.Register)
	}
	if e.Unregister != nil {
		close(e.Unregister)
	}
	if e.ClientMessage != nil {
		close(e.ClientMessage)
	}
}
