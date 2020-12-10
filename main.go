package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	name       int
	connection *websocket.Conn
}

type ChatRoom struct {
	clients []*Client
}

func (chatroom *ChatRoom) Join(name int, conn *websocket.Conn) {
	client := NewClient(rand.Intn(100), conn)
	chatroom.clients = append(chatroom.clients, client)
}

func NewChatRoom() *ChatRoom {
	chatroom := &ChatRoom{
		clients: make([]*Client, 0),
	}
	return chatroom
}

func NewClient(name int, conn *websocket.Conn) *Client {
	client := &Client{
		name:       name,
		connection: conn,
	}

	return client
}

func wsServer(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected...")

	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	chatroom := NewChatRoom()
	chatroom.Join(rand.Intn(100), ws)

	client := NewClient(rand.Intn(100), ws)
	fmt.Println(client)
	client.Reader(chatroom)

}

func (client *Client) Reader(chatroom *ChatRoom) {
	for {
		_, p, err := client.connection.ReadMessage()
		if err != nil {
			fmt.Println("Error11")
			log.Println(err)
			return
		}
		fmt.Printf("user_%d:%s \n", client.name, string(p))
		for _, c := range chatroom.clients {
			c.Writer(string(p), client.name)
		}
	}
}

func (client *Client) Writer(msg string, writer int) {
	wn := fmt.Sprintf("user_%d : %s", writer, msg)
	client.connection.WriteMessage(1, []byte(wn))
}

func router() {
	http.HandleFunc("/", wsServer)
}

func main() {
	fmt.Println("Start...")
	router()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
