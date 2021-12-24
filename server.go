package nyusocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	Subprotocols: []string{"soap"},
}

type Socket struct {
	ctx     context.Context
	events  *Events
	options Options
	hub     *Hub
	start   bool
}

func NewServer(ctx context.Context, options Options) Socket {
	s := Socket{}
	s.ctx = ctx
	s.events = NewEvents()
	s.options = options
	return s
}

func (s Socket) GetEvents() *Events {
	return s.events
}

// Start the socket server
func (s *Socket) Start() error {
	s.start = true
	defer func() {
		s.events.Close()
		s.start = false
	}()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	s.hub = newHub()
	go s.hub.run(s.events)

	r := mux.NewRouter()
	r.HandleFunc("/", index(s.hub, s.events))

	var sErr error
	server := http.Server{Addr: s.options.Addr, Handler: r}
	go func() {
		<-s.ctx.Done()
		s.events.Close()
		if err := server.Shutdown(context.Background()); err != nil {
			sErr = err
		}
	}()

	log.Printf("Server websocket running on %s", s.options.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sErr = err
	}
	log.Println("Server websocket closed")

	return sErr
}

func (s Socket) IsStarting() bool {
	return s.start
}

func (s Socket) GetHub() *Hub {
	return s.hub
}

func index(clientHub *Hub, events *Events) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &Client{hub: clientHub, hash: ksuid.New().String(), Store: NewStore(), path: r.URL.Path, query: r.URL.Query()}

		isAuthorizate := true
		if events.BeforeUpgrade != nil {
			authorize := make(chan bool)
			events.BeforeUpgrade <- BeforeUpgrade{
				Client:    client,
				Authorize: authorize,
			}
			isAuthorizate = <-authorize
			close(authorize)
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !isAuthorizate {
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			t, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			t.Write(NewMessage("Error", "Not authorize", "not_authorize").Send())
			if err := t.Close(); err != nil {
				return
			}
			log.Println("Disconnected user during login (Not authorize)")
			conn.Close()
			return
		}

		client.conn = conn
		client.send = make(chan []byte, 256)
		client.hub.register <- client

		go client.readPump()
		go client.writePump()
	}
}
