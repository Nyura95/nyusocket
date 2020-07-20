package nyusocket

// Events ...
type Events struct {
	Authorization chan Authorization
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
	e.CreateAuthorizationEvent()
	e.CreateClientMessageEvent()
	e.CreateRegisterEvent()
	e.CreateUnregisterEvent()
}

// CreateAuthorizationEvent ...
func (e *Events) CreateAuthorizationEvent() {
	e.Authorization = make(chan Authorization)
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
	if e.Authorization != nil {
		close(e.Authorization)
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
