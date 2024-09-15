package user

import (
	"github.com/gin-gonic/gin"
	"go_web/app/res"
)

func Test1(c *gin.Context) {
	res.Json(c)
	return
}

func Test(c *gin.Context) {
	res.Json(c, res.Data(gin.H{
		"list": []map[string]string{
			{
				"title": "xxxxxx1",
			},
			{
				"title": "xxxxxx2",
			},
			{
				"title": "xxxxxx3",
			},
			{
				"title": "xxxxxx4",
			},
			{
				"title": "xxxxxx5",
			},
		},
	}))
	return
}
