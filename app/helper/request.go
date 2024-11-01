package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetToken(c *gin.Context) (string, error) {
	//SecWebSocketProtocol := c.Request.Header.Get("Sec-WebSocket-Protocol")
	//if SecWebSocketProtocol != "" {
	//	//fmt.Println(SecWebSocketProtocol)
	//}
	token_get := c.DefaultQuery("token", "")
	if token_get != "" {
		return token_get, nil
	}

	token_post := c.DefaultPostForm("token", "")
	if token_post != "" {
		return token_post, nil
	}

	Authorization := c.Request.Header.Get("Authorization")
	if Authorization == "" {
		return "", errors.New("Token 不存在")
	}
	index := strings.Index(Authorization, "Bearer ")
	if index == -1 {
		return "", errors.New("Token Bearer 错误")
	}
	return Authorization[index+len("Bearer "):], nil
}
