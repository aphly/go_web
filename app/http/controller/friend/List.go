package friend

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/http/model"
	"go_web/app/http/service"
	"go_web/app/res"
	"gorm.io/gorm"
)

func List(c *gin.Context) {
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	var friends []model.Friend
	app.DbW().Where("uuid=?", me.Uuid).Where("status=1").
		Preload("FriendUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("nickname", "uuid", "avatar", "remote")
		}).Find(&friends)
	var list []any
	for _, v := range friends {
		v.AvatarPath = service.UploadPath(v.FriendUser.Avatar, v.FriendUser.Remote)
		list = append(list, v)
	}
	res.Json(c, res.Data(gin.H{
		"list": list,
	}))
	return
}
