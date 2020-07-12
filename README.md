# NyuSocket

Basic usage

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"

func main() {
  Events := nyusocket.NewEvents()
	go nyusocket.Start(Events)

	for {
		select {
		case auth := <-Events.Authorization:
			auth.Authorize <- true
		case client := <-Events.Register:
			client.Send <- nyusocket.NewMessage("register", "Hello", "new_register").Send()
			for client := range client.Hub.Clients {
				client.Send <- nyusocket.NewMessage("register", "New client", "new_register").Send()
			}
			log.Printf("New client register alive now : %d", nyusocket.Infos.NbAlive())
		case clients := <-Events.Unregister:
			for client := range clients {
				client.Send <- nyusocket.NewMessage("unregister", "New unregister", "new_unregister").Send()
			}
			log.Printf("Client unregister alive now : %d", nyusocket.Infos.NbAlive())
		}
	}
}
```
