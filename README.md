# NyuSocket

```shell
go get github.com/Nyura95/nyusocket
```

Basic usage

```go
package main

import "fmt"
import "os"
import "github.com/Nyura95/nyusocket"

func main() {
  Events := socket.NewEvents()
	go socket.Start(Events, socket.Options{Port: port})

	for {
		select {
		case auth := <-Events.Authorization:
			// only one hash
			auth.Authorize <- !socket.Infos.Alive(auth.Hash)
		case client := <-Events.Register:
			client.Send <- socket.NewMessage("register", "Hello", "new_register").Send()
			for _, other := range client.Hub.GetOtherClient(client) {
				other.Send <- socket.NewMessage("register", "New client", "new_register").Send()
			}
			log.Printf("New client register alive now : %d", socket.Infos.NbAlive())
		case unregister := <-Events.Unregister:
			for _, client := range unregister.Client.Hub.GetOtherClient(unregister.Client) {
				client.Send <- socket.NewMessage("unregister", "New unregister", "new_unregister").Send()
			}
      log.Printf("Client unregister alive now : %d", socket.Infos.NbAlive())
			unregister.Continue <- true // Mandatory !
		case clientMessage := <-Events.ClientMessage:
			for _, other := range clientMessage.Client.Hub.GetOtherClient(clientMessage.Client) {
				other.Send <- socket.NewMessage("message", clientMessage.Message, "message").Send()
			}
		}
	}
}
```

Exemple client

```tsx
import React, {
  FunctionComponent,
  useCallback,
  useEffect,
  useState,
} from "react";
import useWebSocket from "react-use-websocket";

interface IMessage {
  Action: string;
  Message: string;
  Key: string;
  Created: string;
}

const page: FunctionComponent = () => {
  const [messages, setMessages] = useState<string[]>([]);
  const { readyState, lastMessage, sendMessage } = useWebSocket(
    "ws://localhost:3001"
  );

  // const { readyState, lastMessage, sendMessage } = useWebSocket(
  //   "ws://localhost:3001/anyTokenAuthorization"
  // );

  const triggerMessage = useCallback(() => {
    sendMessage("Hi !!");
  }, []);

  useEffect(() => {
    if (lastMessage) {
      try {
        for (const message of lastMessage.data.split("\n")) {
          const json: IMessage = JSON.parse(message);
          if (
            messages.length < 1 ||
            json.Created !== messages[messages.length - 1].Created
          ) {
            setMessages([...messages, json]);
          }
        }
      } catch (err) {
        console.log(err);
      }
    }
  }, [lastMessage, messages, setMessages]);

  return (
    <Fragment>
      <div>state: {readyState}</div>
      {messages.map((x, index) => (
        <div key={index}>{x.Message}</div>
      ))}
      <button onClick={triggerMessage}>Send Hi!</button>
    </Fragment>
  );
};
```
