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
	client := NewClient(rand.Intn(100), ws)
	fmt.Println(client)
	client.Reader()
}

func (client *Client) Reader() {
	for {
		_, p, err := client.connection.ReadMessage()
		if err != nil {
			fmt.Println("Error11")
			log.Println(err)
			return
		}
		fmt.Printf("user --> %d:%s \n", client.name, string(p))
	}
}

func writer() {

}

func router() {
	http.HandleFunc("/", wsServer)
}

func main() {
	fmt.Println("Start...")
	router()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
