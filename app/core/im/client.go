package im

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go_web/app"
	"go_web/app/core"
	"go_web/app/http/model"
	"net/http"
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

var Ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
	Uuid core.Int64
}

func (this *Client) dispatch(msg []byte) error {
	var data map[string]any
	err := json.Unmarshal(msg, &data)
	if err != nil {
		return err
	}
	val, ok := data["act"]
	if !ok {
		return errors.New("数据格式错误")
	}
	if val == "send" {
		var chatMessage model.ChatMessage
		valData, ok := data["data"]
		if !ok {
			return errors.New("数据格式错误")
		}
		jsonData, err := json.Marshal(valData)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(jsonData, &chatMessage)
		if err != nil {
			return err
		}

		friendClient, ok := this.Hub.Clients[chatMessage.ToId]
		if ok {
			friendClient.Send <- []byte(chatMessage.Message)
		}
		app.DbW().Create(&chatMessage)
	}
	return nil
}

func (this *Client) Read() {
	defer func() {
		this.Hub.Unregister <- this
		this.Conn.Close()
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
		}
	}()
	this.Conn.SetReadLimit(maxMessageSize)
	this.Conn.SetReadDeadline(time.Now().Add(pongWait))
	this.Conn.SetPongHandler(func(string) error { this.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := this.Conn.ReadMessage()
		if err != nil {
			break
		}
		err = this.dispatch(message)
		if err != nil {
			break
		}
	}
}

func (this *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		this.Conn.Close()
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
		}
	}()
	for {
		select {
		case message, ok := <-this.Send:
			this.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				this.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := this.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(this.Send)
			for i := 0; i < n; i++ {
				message_next := <-this.Send
				w.Write(message_next)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			this.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := this.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
