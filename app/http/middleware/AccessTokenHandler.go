package middleware

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

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization := c.Request.Header.Get("Authorization")
		if Authorization == "" {
			res.Json(c, res.Code(1), res.Msg("Token 不存在"))
			c.Abort()
			return
		}
		index := strings.Index(Authorization, "Bearer ")
		if index == -1 {
			res.Json(c, res.Code(1), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		token := Authorization[index+len("Bearer "):]
		de, err := crypt.AesDe(token)
		if err != nil {
			res.Json(c, res.Code(2), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		uuid_token := strings.Split(de, "_")
		parseInt, err := strconv.ParseInt(uuid_token[0], 10, 64)
		if err != nil {
			res.Json(c, res.Code(3), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		user := model.User{
			Uuid: parseInt,
		}
		app.DbW().Where(&user).Take(&user)
		if uuid_token[1] != user.AccessToken {
			res.Json(c, res.Code(4), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		if user.AccessTokenExpire < time.Now().Unix() {
			res.Json(c, res.Code(401), res.Msg("Access Token 过期"))
			c.Abort()
			return
		}
		c.Next()
	}
}
