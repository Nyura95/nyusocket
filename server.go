package nyusocket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

// Start the socket server
func Start(events *Events, options Options) {

	var ShadowLands = newHub()
	go ShadowLands.run(events)

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(ShadowLands, ksuid.New().String(), w, r)
	})

	r.HandleFunc("/{hash}", func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		if events.Authorization != nil {
			authorize := make(chan bool)
			events.Authorization <- Authorization{
				Hash:      hash,
				Authorize: authorize,
			}
			isAuthorizate := <-authorize
			close(authorize)
			if !isAuthorizate {
				closeServeWs("Not authorize", "not_authorize", w, r)
				return
			}
		}
		serveWs(ShadowLands, hash, w, r)
	})

	if err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(options.Port), r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
