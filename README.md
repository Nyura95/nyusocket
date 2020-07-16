# NyuSocket

```shell
go get github.com/Nyura95/nyusocket
```

## Next step of development

- Change the creation of events to make it more 'user frendly'
- Adding log server

## Basic usage

### Sending Message

Server

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	go socket.Start(events, socket.Options{Port: 3000})
	for {
		select {
		case clientMessage := <-events.ClientMessage:
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send <- socket.NewMessage("message", clientMessage.Message, "message").Send()
			}
		}
	}
}
```

### Event before a login user

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"

type storeClient struct {
	token string
}

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	events.CreateAuthorizationEvent()
	go socket.Start(events, socket.Options{Port: 3000})
	for {
		select {
		case authorization := <-events.Authorization:
			// authorize only one 'token' (check the client for pass the query)
			authorization.Client.Store = storeClient{
				token: authorization.Client.Query["token"][0],
			}
			authorization.Client.Hash = authorization.Client.Query["token"][0]
			authorization.Authorize <- !socket.Infos.Alive(authorization.Client)
		case clientMessage := <-events.ClientMessage:
			storeClient := clientMessage.Client.Store.(storeClient)
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send <- socket.NewMessage("message", fmt.Sprintf("%s: %s", storeClient.token, clientMessage.Message), "message").Send()
			}
		}
	}
}
```

### Event on the new client login

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	events.CreateRegisterEvent()
	go socket.Start(events, socket.Options{Port: 3000})
	for {
		select {
		case clientMessage := <-events.ClientMessage:
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send <- socket.NewMessage("message", clientMessage.Message, "message").Send()
			}
		case client := <-events.Register:
			client.Send <- socket.NewMessage("register", "Hello there!", "new_register").Send()
			for _, other := range client.GetOthersClients() {
				other.Send <- socket.NewMessage("register", "New client", "new_register").Send()
			}
		}
	}
}
```

### Event on the client logout

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	events.CreateUnregisterEvent()
	go socket.Start(events, socket.Options{Port: 3000})
	for {
		select {
		case clientMessage := <-events.ClientMessage:
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send <- socket.NewMessage("message", clientMessage.Message, "message").Send()
			}
    case unregister := <-events.Unregister:
      // unregister.Store
			for _, other := range unregister.Hub.GetClients() {
				other.Send <- socket.NewMessage("unregister", "Client unregister", "new_unregister").Send()
			}
		}
	}
}
```
