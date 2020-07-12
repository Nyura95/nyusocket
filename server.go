package nyusocket

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// ShadowLands is a socket hub
var ShadowLands = newHub()

// Start the socket server
func Start(events *Events) {

	r := mux.NewRouter()

	go ShadowLands.run(events)

	r.HandleFunc("/{hash}", func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		if events.Authorization != nil {
			authorize := make(chan bool)
			events.Authorization <- Authorization{Hash: hash, Authorize: authorize}

			isAuthorizate := <-authorize
			close(authorize)
			if !isAuthorizate {
				closeServeWs("Not authorize", "not_authorize", w, r)
				return
			}
		}
		serveWs(ShadowLands, hash, w, r)
	})

	log.Println("Start Socket server on localhost:" + os.Getenv("SOCKET"))
	if err := http.ListenAndServe("127.0.0.1:"+os.Getenv("SOCKET"), r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
