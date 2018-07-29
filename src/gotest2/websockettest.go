// gotest2 project gotest2.go
package gotest2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

const T_FORMAT = "2006-01-02 15:04:05.000"

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Response  string `json:"response,omitempty"`
}

var manager = ClientManager{
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func StartListening() {
	go manager.start()
	http.HandleFunc("/ws", wsPage)

	http.ListenAndServe(":12345", nil)
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			manager.send(jsonMessage, conn)
			fmt.Println("/A new socket has connected.")
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				manager.send(jsonMessage, conn)
			}
			fmt.Println("/A socket has disconnected.")
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			fmt.Println("send notify message.")
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message), Response: string(time.Now().Format(T_FORMAT))})

		fmt.Println("read message:", string(jsonMessage))
		c.send <- jsonMessage

	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				fmt.Println("write close message.")
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			fmt.Println("write message.", string(message))
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("receive request from client")
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	uidrand, _ := uuid.NewV4()
	fmt.Println("client uid = ", uidrand)

	client := &Client{id: uidrand.String(), socket: conn, send: make(chan []byte, 1000)}

	manager.register <- client

	go client.read()
	go client.write()
}
