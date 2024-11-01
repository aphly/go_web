package middleware

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		defer func() {
			if err := recover(); err != nil {
				var obj gin.H
				switch v := err.(type) {
				case error:
					obj = gin.H{
						"code": 501,
						"msg":  v.Error(),
					}
				case gin.H:
					obj = v
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

				endTime := time.Now()
				path := c.Request.URL.Path
				raw := c.Request.URL.RawQuery
				if raw != "" {
					path = path + "?" + raw
				}
				param := []string{
					endTime.Format("2006/01/02 - 15:04:05"),
					strconv.Itoa(obj["code"].(int)),
					endTime.Sub(startTime).String(),
					c.ClientIP(),
					c.Request.Method,
					c.Request.Header.Get("Content-Type"),
					strconv.Itoa(c.Writer.Size()),
					path,
					obj["msg"].(string),
				}
				strLog := strings.Join(param, " | ")
				app.Log().Request(strLog)
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
