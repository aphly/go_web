package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_web/app"
	"go_web/app/core/im"
	"go_web/app/http/model"
)

func Message(c *gin.Context) {
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	ws, err := im.Ws.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		return
	}
	//uuid := c.DefaultQuery("uuid", "")
	//if uuid == "" {
	//	_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "用户错误1"))
	//	return
	//}
	//var me = model.User{}
	//app.DbW().Where("uuid=?", uuid).Take(&me)

	friend_uuid := c.DefaultQuery("friend_uuid", "")
	if friend_uuid == "" {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "用户错误2"))
		return
	}
	var friend model.Friend
	result := app.DbW().Where("uuid=? and status=1", me.Uuid).Where("friend_uuid=?", friend_uuid).Take(&friend)
	if result.Error != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "不是好友关系"))
		return
	}
	meClient := &im.Client{
		Hub:  im.NewHub,
		Conn: ws,
		Uuid: me.Uuid,
		Send: make(chan []byte),
	}
	meClient.Hub.Register <- meClient
	go meClient.Read()
	go meClient.Write()
}
