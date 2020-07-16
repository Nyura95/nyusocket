package nyusocket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Start the socket server
func Start(events *Events, options Options) {

	var hub = newHub()
	go hub.run(events)

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
		serveWs(hub, client, w, r)
	})

	if err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(options.Port), r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
