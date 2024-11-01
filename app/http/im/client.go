package im

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_web/app"
	"go_web/app/core"
	"go_web/app/helper"
	"go_web/app/http/model"
	"gorm.io/gorm"
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

	MessagePagesize = 10
)

var Ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	Uuid   core.Int64
	ChatId core.Int64
	ToId   core.Int64
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
	err := this.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		_ = this.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "pongWait"))
		return
	}
	this.Conn.SetPongHandler(func(string) error {
		this.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := this.Conn.ReadMessage()
		if err != nil {
			_ = this.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
			break
		}
		err = this.Dispatch(message)
		if err != nil {
			_ = this.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
			break
		}
	}
}

func (this *Client) Dispatch(msg []byte) error {
	var data map[string]any
	err := json.Unmarshal(msg, &data)
	if err != nil {
		return err
	}
	val, ok := data["act"]
	if !ok {
		return errors.New("数据格式错误")
	}
	valData, ok := data["data"]
	if !ok {
		return errors.New("数据格式错误")
	}
	DataJson, err := json.Marshal(valData)
	if err != nil {
		return errors.New("数据格式错误")
	}
	if val == "send" {
		var chatMessage model.ChatMessage
		chatMessage.ChatId = this.ChatId
		chatMessage.Uuid = this.Uuid
		chatMessage.ToId = this.ToId
		err = json.Unmarshal(DataJson, &chatMessage)
		if err != nil {
			return err
		}
		err := app.DbW().Create(&chatMessage).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		friendClient, ok := this.Hub.Clients[chatMessage.ToId]
		if ok {
			marshal, _ := json.Marshal(gin.H{
				"msg_type": "res_send",
				"data":     chatMessage,
			})
			friendClient.Send <- marshal
		}
		marshal, err := json.Marshal(gin.H{
			"msg_type": "res_send",
			"data":     chatMessage,
		})
		if err != nil {
			return err
		}
		this.Send <- marshal
	} else if val == "list" {
		var chatMessageQuery model.ChatMessageQuery
		err = json.Unmarshal(DataJson, &chatMessageQuery)
		if err != nil {
			return err
		}
		offset := helper.PageOffset(chatMessageQuery.Page, MessagePagesize)
		list, err := this.ChatMessageList(offset, MessagePagesize, chatMessageQuery.Timestamp)
		if err != nil {
			return err
		}
		this.Send <- list
	}
	return nil
}

func (this *Client) ChatMessageList(offset, pageSize int, timestamp int64) ([]byte, error) {
	var chatMessage []model.ChatMessage
	var count int64 = 0
	if timestamp == 0 {
		app.DbW().Model(&model.ChatMessage{}).Where("chat_id=?", this.ChatId).Count(&count)
		if count > 0 {
			app.DbW().Where("chat_id=?", this.ChatId).Preload("ChatMessageUser", func(db *gorm.DB) *gorm.DB {
				return db.Select("nickname", "uuid")
			}).Order("created_at desc").Offset(offset).Limit(pageSize).Find(&chatMessage)
		}
	} else {
		app.DbW().Where("chat_id=? and created_at<=?", this.ChatId, timestamp).Preload("ChatMessageUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("nickname", "uuid")
		}).Order("created_at desc").Offset(offset).Limit(pageSize).Find(&chatMessage)
	}
	chatMessage_new := make([]model.ChatMessage, len(chatMessage))
	for i := range chatMessage {
		chatMessage_new[len(chatMessage)-1-i] = chatMessage[i]
	}
	var res gin.H
	if timestamp == 0 {
		res = gin.H{
			"msg_type":  "res_list",
			"data":      chatMessage_new,
			"timestamp": time.Now().Unix(),
			"pages":     helper.Pages(int(count), pageSize),
		}
	} else {
		res = gin.H{
			"msg_type": "res_list",
			"data":     chatMessage_new,
		}
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return marshal, nil
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
