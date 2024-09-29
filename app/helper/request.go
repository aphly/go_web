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
	token := c.DefaultQuery("token", "")
	if token != "" {
		return token, nil
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
