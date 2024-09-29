package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_web/app"
	"go_web/app/core/crypt"
	"go_web/app/helper"
	"go_web/app/http/model"
	"go_web/app/res"
	"time"
)

func Register(c *gin.Context) {
	userAuthForm := model.UserAuthForm{}
	err := c.ShouldBind(&userAuthForm)
	if err != nil {
		res.Json(c, res.Code(1), res.Msg(userAuthForm.GetError(err.(validator.ValidationErrors))))
		return
	}
	userAuth := model.UserAuth{}
	userAuth.Id = userAuthForm.Id
	userAuth.IdType = "mobile"
	app.DbW().Where(&userAuth).Take(&userAuth)
	if userAuth.Uuid != 0 {
		res.Json(c, res.Code(1), res.Msg("手机号码已存在"))
		return
	}
	userAuth.Password = crypt.ShaEn(userAuthForm.Password)
	userAuth.LastIp = c.ClientIP()
	userAuth.LastTime = time.Now().Unix()
	userAuth.UserAgent = c.Request.Header.Get("User-Agent")
	userAuth.AcceptLanguage = c.Request.Header.Get("Accept-Language")
	userAuth.Uuid = helper.NewSnowflake.NextID()
	result := app.DbW().Create(&userAuth)
	if result.Error != nil {
		res.Json(c, res.Code(1), res.Msg("错误"))
		return
	}
	user := &model.User{}
	err = user.Add(userAuth.Uuid)
	if err != nil {
		res.Json(c, res.Code(1), res.Msg("错误1"))
		return
	}
	res.Json(c, res.Data(gin.H{
		"user": gin.H{
			"id_type":       userAuth.IdType,
			"id":            userAuth.Id,
			"uuid":          userAuth.Uuid,
			"access_token":  user.EnToken(user.AccessToken),
			"refresh_token": user.EnToken(user.RefreshToken),
			"nickname":      user.Nickname,
		},
	}))
	return
}
