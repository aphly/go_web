package user

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/core/crypt"
	"go_web/app/http/model"
	"go_web/app/res"
	"strconv"
	"strings"
	"time"
)

func RefreshToken(c *gin.Context) {
	Authorization := c.Request.Header.Get("Authorization")
	if Authorization == "" {
		res.Json(c, res.Code(1), res.Msg("Token 不存在"))
		return
	}
	index := strings.Index(Authorization, "Bearer ")
	if index == -1 {
		res.Json(c, res.Code(1), res.Msg("Token 错误"))
		return
	}
	token := Authorization[index+len("Bearer "):]
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
		Uuid: parseInt,
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
			"uuid":         strconv.FormatInt(user.Uuid, 10),
			"access_token": user.EnToken(user.AccessToken),
		},
	}))
	return
}
