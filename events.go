package nyusocket

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
func (e *Events) CreateBeforeUpgradeEvent() <-chan BeforeUpgrade {
	e.BeforeUpgrade = make(chan BeforeUpgrade)
	return e.BeforeUpgrade
}

// CreateRegisterEvent ...
func (e *Events) CreateRegisterEvent() <-chan *Client {
	e.Register = make(chan *Client)
	return e.Register
}

// CreateUnregisterEvent ...
func (e *Events) CreateUnregisterEvent() <-chan Unregister {
	e.Unregister = make(chan Unregister)
	return e.Unregister
}

// CreateClientMessageEvent ...
func (e *Events) CreateClientMessageEvent() chan ClientMessage {
	e.ClientMessage = make(chan ClientMessage)
	return e.ClientMessage
}

// Close all chan event
func (e *Events) Close() {
	if e.BeforeUpgrade != nil {
		close(e.BeforeUpgrade)
		e.BeforeUpgrade = nil
	}
	if e.Register != nil {
		close(e.Register)
		e.Register = nil
	}
	if e.Unregister != nil {
		close(e.Unregister)
		e.Unregister = nil
	}
	if e.ClientMessage != nil {
		close(e.ClientMessage)
		e.ClientMessage = nil
	}
}
