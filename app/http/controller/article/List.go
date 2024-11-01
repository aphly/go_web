package article

import (
	"github.com/gin-gonic/gin"
	"go_web/app/res"
)

var a1 = make(map[string]string)

func List(c *gin.Context) {
	key := c.Query("key")
	val := c.Query("val")
	a1[key] = val
	res.Json(c, res.Data(gin.H{
		"list": []map[string]string{
			{
				"title": "xxxxxx1",
			},
			{
				"title": "xxxxxx2",
			},
		},
		"map": a1,
	}))
	return
}
