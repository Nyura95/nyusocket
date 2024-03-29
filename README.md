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
	socket := nyusocket.NewServer(nyusocket.Options{Addr: "127.0.0.1:3001"})
	events := socket.GetEvents()

	ctx, cancel := context.WithCancel(context.Background())
	go socket.Start(ctx)
	defer cancel()

	events.CreateClientMessageEvent()
	for {
		select {
		case clientMessage, alive := <-events.ClientMessage:
			if !alive {
				break
			}
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send(nyusocket.NewMessage("message", clientMessage.Message, "message").Send())
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
	socket := nyusocket.NewServer(nyusocket.Options{Addr: "127.0.0.1:3001"})
	events := socket.GetEvents()

	ctx, cancel := context.WithCancel(context.Background())
	go socket.Start(ctx)
	defer cancel()

	events.CreateClientMessageEvent()
	events.CreateAuthorizationEvent()
	for {
		select {
		case authorization, alive := <-events.Authorization:
			if !alive {
				break
			}
			// authorize only one 'token' (check the client for pass the query)
			authorization.Client.Store = storeClient{
				token: authorization.Client.Query["token"][0],
			}
			authorization.Client.Hash = authorization.Client.Query["token"][0]
			authorization.Authorize <- !nyusocket.Infos.Alive(authorization.Client) // chan authorization.Authorize return an boolean, if false the client is unregister
		case clientMessage, alive := <-events.ClientMessage:
			if !alive {
				break
			}
			storeClient := clientMessage.Client.Store.(storeClient) // get store client
			for _, other := range clientMessage.Client.GetOthersClients() { // get all other clients actually registered
				other.Send(nyusocket.NewMessage("message", fmt.Sprintf("%s: %s", storeClient.token, clientMessage.Message), "message").Send()) // send customer's message to others
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
  socket := nyusocket.NewServer(nyusocket.Options{Addr: "127.0.0.1:3001"})
	events := socket.GetEvents()

	ctx, cancel := context.WithCancel(context.Background())
	go socket.Start(ctx)
	defer cancel()

	events.CreateClientMessageEvent()
	events.CreateRegisterEvent()

	for {
		select {
		case clientMessage, alive := <-events.ClientMessage:
			if !alive {
				break
			}
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send(nyusocket.NewMessage("message", clientMessage.Message, "message").Send())
			}
		case client, alive := <-events.Register: // new client registered
			if !alive {
				break
			}
			client.Send(nyusocket.NewMessage("register", "Hello there!", "new_register").Send()) // send a message to the client
			for _, other := range client.GetOthersClients() {
				other.Send(nyusocket.NewMessage("register", "New client", "new_register").Send()) // tell others that a new customer is registered
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
  socket := nyusocket.NewServer(nyusocket.Options{Addr: "127.0.0.1:3001"})
	events := socket.GetEvents()

	ctx, cancel := context.WithCancel(context.Background())
	go socket.Start(ctx)
	defer cancel()

	events.CreateClientMessageEvent()
	events.CreateUnregisterEvent()

	for {
		select {
		case clientMessage, alive := <-events.ClientMessage:
			if !alive {
				break
			}
			for _, other := range clientMessage.Client.GetOthersClients() {
				other.Send(nyusocket.NewMessage("message", clientMessage.Message, "message").Send())
			}
    case unregister, alive := <-events.Unregister:
			if !alive {
				break
			}
      // unregister.Store (you can access to the store client if needed)
			for _, other := range unregister.Hub.GetClients() {
				other.Send(nyusocket.NewMessage("unregister", "Client unregister", "new_unregister").Send()) // tell others that a new customer is logout
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
  socket := nyusocket.NewServer(nyusocket.Options{Addr: "127.0.0.1:3001"})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	socket.Start(ctx)

	log.Println("server closed")
}
```

## Client side

### Message server example

```javascript
{ Action: string, Message: string, Key: string, Created: Date }
```

### Who to use ? (see client.html)

```javascript
var states = {
  [WebSocket.CONNECTING]: "CONNECTING",
  [WebSocket.OPEN]: "OPEN",
  [WebSocket.CLOSING]: "CLOSING",
  [WebSocket.CLOSED]: "CLOSED",
};
var _messages = [];

var client = new WebSocket("ws://...");
client.onopen = function (event) {
  console.log(event);
  document.getElementById("state").innerHTML =
    "state: " + states[client.readyState];
};
client.onmessage = function (msg) {
  console.log("New message: ", msg);
  var messages = msg.data.split("\n");
  for (let i = 0; i < messages.length; i++) {
    var message = messages[i];
    try {
      message = JSON.parse(message);
    } catch (err) {
      console.log(err);
    }
    _messages.push(message);
    var tchat = document.getElementById("messages");
    tchat.value = _messages
      .map(function (x) {
        return "Action: " + x.Action + " | " + x.Message;
      })
      .join("\r\n");
    tchat.scrollTo({ top: tchat.scrollHeight });
  }
};
```
