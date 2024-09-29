package middleware

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

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helper.GetToken(c)
		if err != nil {
			res.Json(c, res.Code(401), res.Msg(err.Error()))
			c.Abort()
			return
		}
		de, err := crypt.AesDe(token)
		if err != nil {
			res.Json(c, res.Code(401), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		uuid_token := strings.Split(de, "_")
		uuidInt, err := strconv.ParseInt(uuid_token[0], 10, 64)
		if err != nil {
			res.Json(c, res.Code(401), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		user := model.User{
			Uuid: core.Int64(uuidInt),
		}
		app.DbW().Where(&user).Take(&user)
		if uuid_token[1] != user.AccessToken {
			res.Json(c, res.Code(401), res.Msg("Token 错误"))
			c.Abort()
			return
		}
		if user.AccessTokenExpire < time.Now().Unix() {
			res.Json(c, res.Code(402), res.Msg("Access Token 过期"))
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Set("uuid", user.Uuid)
		c.Next()
		//SecWebSocketProtocol := c.Request.Header.Get("Sec-WebSocket-Protocol")
		//if SecWebSocketProtocol != "" {
		//	//c.Header("Sec-WebSocket-Protocol", SecWebSocketProtocol)
		//}

	}
}
