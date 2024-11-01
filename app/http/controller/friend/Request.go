package friend

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_web/app"
	"go_web/app/http/model"
	"go_web/app/res"
	"strconv"
)

func RequestAdd(c *gin.Context) {
	friendRequestForm := model.FriendRequestForm{}
	err := c.ShouldBind(&friendRequestForm)
	if err != nil {
		res.Json(c, res.Code(1), res.Msg(friendRequestForm.GetError(err.(validator.ValidationErrors))))
		return
	}
	userAuth := model.UserAuth{}
	app.DbW().Where("id=?", friendRequestForm.Id).Take(&userAuth)
	if userAuth.Uuid == 0 {
		res.Json(c, res.Code(1), res.Msg("用户不存在"))
		return
	}
	you := model.User{}
	app.DbW().Where("uuid=?", userAuth.Uuid).Take(&you)
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	friendRequest := model.FriendRequest{}
	friendRequest.FromUuid = me.Uuid
	friendRequest.FromName = me.Nickname
	friendRequest.ToUuid = you.Uuid
	friendRequest.ToName = you.Nickname
	app.DbW().Save(&friendRequest)
	res.Json(c, res.Msg("请求成功"))
	return
}

func RequestList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20
	offset := (page - 1) * pageSize
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	var list []model.FriendRequest
	app.DbW().Where("from_uuid=?", me.Uuid).Or("to_uuid=?", me.Uuid).Offset(offset).Limit(pageSize).Find(&list)
	res.Json(c, res.Data(gin.H{
		"list": list,
	}))
	return
}
func RequestOp(c *gin.Context) {
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)
	var friendRequestOpForm model.FriendRequestOpForm
	err := c.ShouldBind(&friendRequestOpForm)
	if err != nil {
		res.Json(c, res.Code(1), res.Msg(friendRequestOpForm.GetError(err.(validator.ValidationErrors))))
		return
	}
	var friendRequest model.FriendRequest
	app.DbW().Where("status=0").
		Where("id=?", friendRequestOpForm.Id).
		Where("to_uuid=? ", me.Uuid).Take(&friendRequest)
	friendRequest.Status = friendRequestOpForm.Status
	result := app.DbW().Save(&friendRequest)
	if result.Error != nil {
		res.Json(c, res.Msg("错误"))
		return
	}
	var friends = []model.Friend{
		{
			Uuid:       me.Uuid,
			FriendUuid: friendRequest.FromUuid,
			Status:     1,
		},
		{
			Uuid:       friendRequest.FromUuid,
			FriendUuid: me.Uuid,
			Status:     1,
		},
	}
	app.DbW().Create(&friends)
	res.Json(c)
	return
}
