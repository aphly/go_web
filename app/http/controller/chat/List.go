package chat

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/http/model"
	"go_web/app/res"
)

func List(c *gin.Context) {
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	var chat []model.Chat
	app.DbW().Where("uuid=?", me.Uuid).Find(&chat)
	res.Json(c, res.Data(gin.H{
		"list": chat,
	}))
	return
}
