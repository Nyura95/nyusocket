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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"soap"},
}

type Server struct {
	events  *Events
	http    http.Server
	options Options
}

func NewServer(options Options) Server {
	s := Server{events: NewEvents(), options: options}

	clientHub := newHub()
	go clientHub.run(s.events)

	r := mux.NewRouter()
	r.HandleFunc("/", index(clientHub, s.events))

	s.http = http.Server{Addr: options.Addr, Handler: r}

	return s
}

func (s Server) Start(ctx context.Context) {

	go func() {
		<-ctx.Done()
		if err := s.http.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		s.events.Close()
	}()

	log.Printf("Server websocket running on %s", s.options.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panicln(err)
	}
	log.Println("Server websocket closed")

}

func (s Server) GetEvents() *Events {
	return s.events
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
