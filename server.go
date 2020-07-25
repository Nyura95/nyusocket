package nyusocket

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var clientHub *Hub

// Start the socket server
func Start(events *Events, options Options) {

	clientHub = newHub()
	go clientHub.run(events)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client := &NewClient{Query: r.URL.Query()}
		if events.Authorization != nil {
			authorize := make(chan bool)
			events.Authorization <- Authorization{
				Client:    client,
				Authorize: authorize,
			}
			isAuthorizate := <-authorize
			close(authorize)
			if !isAuthorizate {
				closeServeWs("Not authorize", "not_authorize", w, r)
				return
			}
		}
		serveWs(clientHub, client, w, r)
	})

	log.Printf("Server websocket running on %s", options.Addr)
	if err := http.ListenAndServe(options.Addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// GetClientHub ...
func GetClientHub() *Hub {
	return clientHub
}
