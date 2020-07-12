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

Exemple client

```tsx
import useWebSocket from "react-use-websocket";

interface IMessage {
  Action: string;
  Message: string;
  Key: string;
  Created: string;
}

const page: FunctionComponent = () => {
  const [messages, setMessages] = useState<string[]>([]);
  const { readyState, lastMessage } = useWebSocket(
    "ws://localhost:3001/anyToken"
  );

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
    </Fragment>
  );
};
```
