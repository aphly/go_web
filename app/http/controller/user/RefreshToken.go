package user

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/core"
	"go_web/app/core/crypt"
	"go_web/app/helper"
	"go_web/app/http/model"
	"go_web/app/res"
	"strconv"
	"strings"
	"time"
)

func RefreshToken(c *gin.Context) {
	token, err := helper.GetToken(c)
	if err != nil {
		res.Json(c, res.Code(1), res.Msg(err.Error()))
		return
	}
	de, err := crypt.AesDe(token)
	if err != nil {
		res.Json(c, res.Code(2), res.Msg("Token 错误"))
		return
	}

	uuid_token := strings.Split(de, "_")
	parseInt, err := strconv.ParseInt(uuid_token[0], 10, 64)
	if err != nil {
		res.Json(c, res.Code(3), res.Msg("Token 错误"))
		return
	}
	user := model.User{
		Uuid: core.Int64(parseInt),
	}
	app.DbW().Where(&user).Take(&user)
	if uuid_token[1] != user.RefreshToken {
		res.Json(c, res.Code(4), res.Msg("Token 错误"))
		return
	}
	now := time.Now().Unix()
	if user.RefreshTokenExpire < now {
		res.Json(c, res.Code(5), res.Msg("Refresh Token 过期"))
		return
	}
	user.GenAccessToken(now)
	app.DbW().Save(&user)
	res.Json(c, res.Data(gin.H{
		"user": gin.H{
			"uuid":         user.Uuid,
			"access_token": user.EnToken(user.AccessToken),
		},
	}))
	return
}
