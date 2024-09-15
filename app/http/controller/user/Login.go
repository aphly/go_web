package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_web/app"
	"go_web/app/core/crypt"
	"go_web/app/http/model"
	"go_web/app/res"
	"strconv"
)

func Login(c *gin.Context) {
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
	if userAuth.Uuid == 0 {
		res.Json(c, res.Code(1), res.Msg("手机号码不存在"))
		return
	}
	if crypt.ShaEn(userAuthForm.Password) != userAuth.Password {
		res.Json(c, res.Code(1), res.Msg("密码错误"))
		return
	}
	user := &model.User{
		Uuid: userAuth.Uuid,
	}
	app.DbW().Where(&user).Take(&user)
	user.GenToken()
	app.DbW().Save(&user)
	res.Json(c, res.Data(gin.H{
		"user": gin.H{
			"id_type":       userAuth.IdType,
			"id":            userAuth.Id,
			"uuid":          strconv.FormatInt(userAuth.Uuid, 10),
			"access_token":  user.EnToken(user.AccessToken),
			"refresh_token": user.EnToken(user.RefreshToken),
			"nickname":      user.Nickname,
		},
	}))
	return
}
