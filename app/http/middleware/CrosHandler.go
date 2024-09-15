package middleware

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"net/http"
)

func CrosHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		for _, domain := range app.Config.Cors.Origin {
			if origin == domain {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "*")
				c.Header("Access-Control-Allow-Headers", "*")
			}
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
