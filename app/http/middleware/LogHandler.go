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
		param := []string{
			endTime.Format("2006/01/02 - 15:04:05"),
			strconv.Itoa(c.Writer.Status()),
			endTime.Sub(startTime).String(),
			c.ClientIP(),
			c.Request.Method,
			c.Request.Header.Get("Content-Type"),
			strconv.Itoa(c.Writer.Size()),
			path,
			c.Errors.ByType(1).String(),
		}

		strLog := strings.Join(param, " | ")
		app.Log().Request(strLog)
	}
}
