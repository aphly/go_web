package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_web/app/core"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域请求
	},
}

type Hub struct {
	Clients    map[core.Int64]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func (this *Hub) Run() {
	for {
		select {
		case client := <-this.Register:
			this.Clients[(*client).Uuid] = client
			msg := []byte("连接成功" + strconv.FormatInt(int64((*client).Uuid), 10))
			client.Conn.WriteMessage(1, msg)
		case client := <-this.Unregister:
			if _, ok := this.Clients[(*client).Uuid]; ok {
				delete(this.Clients, (*client).Uuid)
				close(client.Send)
			}
		}
	}
}

var NewHub = &Hub{
	Clients:    make(map[core.Int64]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan []byte),
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//defer conn.Close()

	u, err := url.Parse(r.URL.String())
	if err != nil {
		panic(err)
	}
	queryParams := u.Query()
	atoi, _ := strconv.Atoi(queryParams.Get("uuid"))
	uuid := core.Int64(atoi)

	//atoi1, _ := strconv.Atoi(queryParams.Get("friend_uuid"))
	//friend_uuid := core.Int64(atoi1)

	meClient := &Client{
		Hub:  NewHub,
		Conn: conn,
		Uuid: uuid,
		Send: make(chan []byte),
	}
	//fmt.Println(meClient)
	meClient.Hub.Register <- meClient
	go meClient.Read()
	go meClient.Write()
}

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
	Uuid core.Int64
}

func (this *Client) Read() {
	for {
		mt, message, err := this.Conn.ReadMessage()
		if err != nil {
			fmt.Println("x1", err)
			break
		}
		fmt.Println(mt, string(message))
		this.Send <- message
		fmt.Println("xxxxx")
	}
}

func (this *Client) Write() {
	fmt.Println("write start")
	for {
		err := this.Conn.WriteMessage(websocket.TextMessage, <-this.Send)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

}

func TestUnit(t *testing.T) {
	go NewHub.Run()
	http.HandleFunc("/message", echo)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
