# NyuSocket

```shell
go get github.com/Nyura95/nyusocket
```

## Next step of development

- Change the creation of events to make it more 'user frendly'

## Basic usage

### Sending Message

Server

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"
import "context"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	go socket.Start(context.Background(), socket.Options{Addr: "127.0.0.1:3000"})
	for {
		select {
		case clientMessage := <-events.ClientMessage:
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send(socket.NewMessage("message", clientMessage.Message, "message").Send())
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
import "context"
import "github.com/Nyura95/nyusocket"

type storeClient struct {
	token string
}

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	events.CreateAuthorizationEvent()
	go socket.Start(context.Background(), events, socket.Options{Addr: "127.0.0.1:3000"})
	for {
		select {
		case authorization := <-events.Authorization:
			// authorize only one 'token' (check the client for pass the query)
			authorization.Client.Store = storeClient{
				token: authorization.Client.Query["token"][0],
			}
			authorization.Client.Hash = authorization.Client.Query["token"][0]
			authorization.Authorize <- !socket.Infos.Alive(authorization.Client) // chan authorization.Authorize return an boolean, if false the client is unregister
		case clientMessage := <-events.ClientMessage:
			storeClient := clientMessage.Client.Store.(storeClient) // get store client
			for _, other := range clientMessage.Client.GetOthersClients() { // get all other clients actually registered
				other.Send(socket.NewMessage("message", fmt.Sprintf("%s: %s", storeClient.token, clientMessage.Message), "message").Send()) // send customer's message to others
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
import "context"
import "github.com/Nyura95/nyusocket"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	events.CreateRegisterEvent()
	go socket.Start(context.Background(), events, socket.Options{Addr: "127.0.0.1:3000"})
	for {
		select {
		case clientMessage := <-events.ClientMessage:
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send(socket.NewMessage("message", clientMessage.Message, "message").Send())
			}
		case client := <-events.Register: // new client registered
			client.Send(socket.NewMessage("register", "Hello there!", "new_register").Send()) // send a message to the client
			for _, other := range client.GetOthersClients() {
				other.Send(socket.NewMessage("register", "New client", "new_register").Send()) // tell others that a new customer is registered
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
import "context"
import "github.com/Nyura95/nyusocket"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()
	events.CreateUnregisterEvent()
	go socket.Start(context.Background(), events, socket.Options{Addr: "127.0.0.1:3000"})
	for {
		select {
		case clientMessage := <-events.ClientMessage:
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send(socket.NewMessage("message", clientMessage.Message, "message").Send())
			}
    case unregister := <-events.Unregister:
      // unregister.Store (you can access to the store client if needed)
			for _, other := range unregister.Hub.GetClients() {
				other.Send(socket.NewMessage("unregister", "Client unregister", "new_unregister").Send()) // tell others that a new customer is logout
			}
		}
	}
}
```

### Stop server

```go
package main

import "fmt"
import "os"
import "time"
import "github.com/Nyura95/nyusocket"
import "context"

func main() {
  events := socket.NewEvents()
	defer events.Close()

	events.CreateClientMessageEvent()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	socket.Start(ctx, socket.Options{Addr: "127.0.0.1:3000"})
	log.Println("server closed")
}
```
