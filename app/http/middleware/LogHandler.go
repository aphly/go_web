package middleware

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"strconv"
	"strings"
	"time"
)

type LogFormatterParams struct {
	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Method       string
	Path         string
	ErrorMessage string
	isTerm       bool
	BodySize     int
	Keys         map[string]any
}

func LogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		endTime := time.Now()
		param := make([]string, 9)
		param[0] = endTime.Format("2006/01/02 - 15:04:05")
		param[1] = strconv.Itoa(c.Writer.Status())
		param[2] = endTime.Sub(startTime).String()
		param[3] = c.ClientIP()
		param[4] = c.Request.Method
		param[5] = strconv.Itoa(c.Writer.Status())
		param[6] = c.Errors.ByType(1).String()
		param[7] = strconv.Itoa(c.Writer.Size())
		param[8] = path

		strLog := strings.Join(param, " | ")
		app.Log().Request(strLog)
	}
}
