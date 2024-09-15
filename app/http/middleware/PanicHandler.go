package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var obj any
				switch v := err.(type) {
				case error:
					obj = gin.H{
						"code": 501,
						"msg":  v.Error(),
					}
				case gin.H:
					obj = err
				case string:
					obj = gin.H{
						"code": 502,
						"msg":  v,
					}
				default:
					obj = gin.H{
						"code": 503,
						"msg":  "error",
					}
				}
				c.JSON(http.StatusOK, obj)
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
