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
