package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var clientMap = map[int]*Client{}

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

func WsServer(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected...")

	err = ws.WriteMessage(1, []byte("Start Chat!"))
	if err != nil {
		log.Println(err)
	}

	rnN := rand.Intn(100)
	client := NewClient(rnN, ws)
	clientMap[rnN] = client
	fmt.Println(clientMap)
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
		fmt.Printf("user_%d:%s \n", client.name, string(p))
		for k, c := range clientMap {
			if k != client.name {
				c.Writer(string(p), client.name)
			}
		}
	}
}

func (client *Client) Writer(msg string, writer int) {
	wn := fmt.Sprintf("user_%d : %s", writer, msg)
	client.connection.WriteMessage(1, []byte(wn))
}

func CreateChatGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func router() {
	http.HandleFunc("/", WsServer)
	http.HandleFunc("/CreateChatGroup", CreateChatGroup)
}

func main() {
	fmt.Println("Start...")
	router()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
