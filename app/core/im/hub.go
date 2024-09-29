package im

import (
	"go_web/app/core"
	"strconv"
)

var NewHub *Hub

func init() {
	NewHub = &Hub{
		Clients:    make(map[core.Int64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
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
		case message := <-this.Broadcast:
			for uuid, client := range this.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(this.Clients, uuid)
				}
			}
		}
	}
}
