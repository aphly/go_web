package chat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_web/app"
	"go_web/app/helper"
	"go_web/app/http/im"
	"go_web/app/http/model"
	"go_web/app/http/service"
	"go_web/app/res"
	"os"
	"strconv"
	"strings"
	"time"
)

func Message(c *gin.Context) {
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	ws, err := im.Ws.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		return
	}

	friend_uuid := c.DefaultQuery("friend_uuid", "")
	if friend_uuid == "" {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "用户错误2"))
		return
	}

	var friend model.Friend
	err = app.DbW().Where("uuid=? and status=1", me.Uuid).Where("friend_uuid=?", friend_uuid).Take(&friend).Error
	if err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "不是好友关系"))
		return
	}
	var MinMax string
	if me.Uuid < friend.FriendUuid {
		MinMax = fmt.Sprintf("%v_%v", me.Uuid, friend.FriendUuid)
	} else {
		MinMax = fmt.Sprintf("%v_%v", friend.FriendUuid, me.Uuid)
	}
	var chat = model.Chat{
		Uuid:        me.Uuid,
		Subject:     "",
		MinMax:      MinMax,
		LastMessage: "",
	}
	app.DbW().Where("min_max=?", MinMax).FirstOrCreate(&chat)

	var meChatMember = model.ChatMember{
		ChatId: chat.Id,
		Uuid:   me.Uuid,
		IsNew:  0,
	}
	app.DbW().Where("chat_id=? and uuid=?", chat.Id, me.Uuid).FirstOrCreate(&meChatMember)

	var friendChatMember = model.ChatMember{
		ChatId: chat.Id,
		Uuid:   friend.FriendUuid,
		IsNew:  0,
	}
	app.DbW().Where("chat_id=? and uuid=?", chat.Id, friend.FriendUuid).FirstOrCreate(&friendChatMember)

	meClient := &im.Client{
		Hub:    im.NewHub,
		Conn:   ws,
		Uuid:   me.Uuid,
		Send:   make(chan []byte),
		ChatId: chat.Id,
		ToId:   friend.FriendUuid,
	}
	meClient.Hub.Register <- meClient
	go meClient.Read()
	go meClient.Write()
	list, err := meClient.ChatMessageList(0, im.MessagePagesize, 0)
	if err != nil {
		return
	}
	meClient.Send <- list
}

func Image(c *gin.Context) {
	file, _ := c.FormFile("file")
	stringArr := file.Header["Content-Type"]
	var upload = &helper.Upload{
		ConutConf: 1,
		SizeConf:  1,
		TypeConf: []any{
			"image/jpeg",
			"image/png",
		},
	}
	if !upload.SizeLimit(file.Size) {
		res.Json(c, res.Code(1), res.Msg("大小超过"+strconv.FormatInt(upload.SizeConf, 10)+"M"))
		return
	}

	if !upload.TypeLimit(stringArr) {
		res.Json(c, res.Code(2), res.Msg("格式只支持"+helper.SliceJoin(upload.TypeConf, "、")))
		return
	}
	var now = time.Now()
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)

	var timeDir = "/public/chat/" + fmt.Sprintf("%d/%d%d/", now.Year(), now.Month(), now.Day()) + helper.UuidPath(me.Uuid) + "/"
	err := os.MkdirAll("."+timeDir, 0755)
	if err != nil {
		res.Json(c, res.Code(3), res.Msg("文件夹错误"))
		return
	}
	arr := strings.Split(file.Filename, ".")
	if len(arr) != 2 {
		res.Json(c, res.Code(4), res.Msg("上传文件格式错误"))
		return
	}
	var savePath = timeDir + helper.RandStr(16) + "." + arr[1]
	var uploadPath = "." + savePath
	err = c.SaveUploadedFile(file, uploadPath)
	if err != nil {
		res.Json(c, res.Code(5), res.Msg("上传错误"))
		return
	}
	res.Json(c, res.Data(gin.H{
		"image": service.UploadPath(savePath, 0),
	}))
	return
}
